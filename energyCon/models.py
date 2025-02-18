import datetime

from pydantic import BaseModel


class Record(BaseModel):
    power1: int
    power2: int
    power3: int
    voltage1: int
    voltage2: int
    voltage3: int
    current1: int
    current2: int
    current3: int
    date: datetime.datetime
