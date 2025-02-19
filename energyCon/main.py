import datetime
from typing import Annotated
from fastapi import FastAPI, status, Form
from models import Record, engine
from sqlmodel import Session, select


app = FastAPI()


@app.post("/rec", status_code=status.HTTP_201_CREATED)
async def new_record(power1: Annotated[str, Form()],
                     power2: Annotated[str, Form()],
                     power3: Annotated[str, Form()],
                     voltage1: Annotated[str, Form()],
                     voltage2: Annotated[str, Form()],
                     voltage3: Annotated[str, Form()],
                     current1: Annotated[str, Form()],
                     current2: Annotated[str, Form()],
                     current3: Annotated[str, Form()],
                     device: Annotated[str, Form()],):
    rec = Record(power1=power1,
                 power2=power2,
                 power3=power3,
                 voltage1=voltage1,
                 voltage2=voltage2,
                 voltage3=voltage3,
                 current1=current1,
                 current2=current2,
                 current3=current3,
                 device=device,
                 date=datetime.datetime.now())
    with Session(engine) as session:
        session.add(rec)
        session.commit()
    return None


@app.get("/")
async def main_page():
    with Session(engine) as session:
        statement = select(Record)
        results = session.exec(statement)
        resp = []
        for r in results:
            resp.append(r)
    return resp
