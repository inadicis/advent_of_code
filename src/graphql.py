import dataclasses
import logging
from datetime import datetime
from functools import cached_property
from random import randint
from typing import Self

from pydantic import BaseModel, Field, field_validator
import strawberry
from strawberry.fastapi import GraphQLRouter
from strawberry.types import Info
from strawberry.types.info import RootValueType

from graphql import NoSchemaIntrospectionCustomRule
from strawberry.extensions import (
    AddValidationRules,
    MaxAliasesLimiter,
    MaxTokensLimiter,
    QueryDepthLimiter,
    SchemaExtension,
)



from src.config.env_settings import get_env_settings


class NoAuthenticationCredentials(Exception):
    ...


@dataclasses.dataclass
class MyExampleModel:
    content: str
    amount: int = 1
    type: str = "A"




@strawberry.interface
class BaseType:
    id: str
    created_at: strawberry.Private[datetime]


@strawberry.type
class MyExampleType(BaseType):
    content: str
    amount: int

    @classmethod
    def from_my_example_model(cls, example: MyExampleModel) -> Self:
        return cls(content=example.content, amount=example.amount, id=randint(1, 1000), created_at=datetime.utcnow())


class DictExportableMixin:
    def dict(self) -> dict:
        return {key: value for key, value in self.__dict__.items() if value is not None}


class NotFound(Exception):
    ...


class QueryOptions(BaseModel):
    limit: int = Field(50, ge=1)
    offset: int = Field(0, ge=0)
    order_by: str | tuple[str, int] | list[tuple[str, int]] = Field(
        [("_id", 1)],
        description=(
            "a list of fields (string name and bool is ascending) to sort by. A simple "
            "string or a unique tuple are also accepted (and transformed to a list of "
            "clauses)"
        ),
    )

    @field_validator("order_by")
    @classmethod
    def cast_order_by_to_list(cls, v) -> list[tuple[str, int]]:
        if isinstance(v, str):
            return [(v, 1)]
        if isinstance(v, tuple):
            return [v]
        return v


async def my_resolver(options: QueryOptions) -> list[MyExampleModel]:
    # TODO-CONFIG an actual resolver would read Database or similar
    return [MyExampleModel(content="fake_content", amount=i) for i in range(10)]
    # return (
    #     await Story.find()
    #     .limit(options.limit)
    #     .skip(options.offset)
    #     .sort(options.order_by)
    #     .to_list()
    # )


@strawberry.input
class QueryOptionsInput(DictExportableMixin):
    limit: int | None = None
    offset: int | None = None
    order_by: list[tuple[str, bool]] | None = None


@strawberry.type
class Query:
    @strawberry.field
    async def examples(
        self,
        options: QueryOptionsInput,
        info: Info,
    ) -> list[MyExampleType]:
        
        return [
            MyExampleType.from_my_example_model(s)
            for s in await my_resolver(options=QueryOptions(**options.dict()))
        ]


# @strawberry.type
# class Mutation:
#     @strawberry.mutation
#     def add_transaction(self, transaction_in: TransactionIn, info: Info) -> TransactionType:
#         authorize_user(
#             user=info.context.auth_user,
#             required_permissions={Permissions.WRITE_OWN_NAMESPACE},
#         )
#         options = TransactionCreationOptions(**transaction_in.dict())
#         transaction = transaction_creation_resolver(
#             options=options,
#             db_session=info.context.db_session,
#             user_id=info.context.auth_user.id,
#         )
#
#         return TransactionType.from_transaction(transaction)


schema = strawberry.Schema(
    query=Query,
    # mutation=Mutation,
)

graphql_app = GraphQLRouter(
    schema,
    
)


class LogErrorStackTrace(SchemaExtension):
    def on_execute(self):
        try:
            yield
        except Exception as e:
            logging.exception(e)
            raise


schema_extensions = [
    QueryDepthLimiter(max_depth=10),
    MaxAliasesLimiter(max_alias_count=15),
    MaxTokensLimiter(max_token_count=1000),
]
if get_env_settings().debug:
    schema_extensions.append(AddValidationRules([NoSchemaIntrospectionCustomRule]))
    schema_extensions.append(LogErrorStackTrace)

schema = strawberry.Schema(
    query=Query,

    extensions=schema_extensions,
)

graphql_app = GraphQLRouter(
    schema,
    
    graphiql=get_env_settings().debug,
)
