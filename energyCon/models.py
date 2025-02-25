import datetime

from sqlmodel import Field, SQLModel, create_engine


class Record(SQLModel, table=True):
    id: int | None = Field(default=None, primary_key=True)
    power1: int
    power2: int
    power3: int
    voltage1: int
    voltage2: int
    voltage3: int
    current1: int
    current2: int
    current3: int
    device: int | None = Field(default=1)
    date: datetime.datetime


sqlite_file_name = "database.db"
sqlite_url = f"sqlite:///{sqlite_file_name}"
engine = create_engine(sqlite_url)

SQLModel.metadata.create_all(engine)
