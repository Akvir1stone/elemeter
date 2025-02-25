import datetime
from typing import Annotated

import fastapi.responses
from fastapi.encoders import jsonable_encoder
from fastapi import FastAPI, status, Form, Request, staticfiles
from fastapi.templating import Jinja2Templates
from models import Record, engine
from sqlmodel import Session, select


app = FastAPI()
app.mount("/static", staticfiles.StaticFiles(directory="static"), name="static")
templates = Jinja2Templates(directory="templates")


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
                     device: Annotated[int, Form()],):
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


@app.post("/no_connection", status_code=status.HTTP_200_OK)
async def no_connection(device: Annotated[int, Form()],):
    rec = Record(power1=0,
                 power2=0,
                 power3=0,
                 voltage1=0,
                 voltage2=0,
                 voltage3=0,
                 current1=0,
                 current2=0,
                 current3=0,
                 device=device,
                 date=datetime.datetime.now())
    with Session(engine) as session:
        session.add(rec)
        session.commit()
    return None


@app.get("/")
async def main_page(request: Request, device: str | None = "1", date: datetime.date | None = None):
    with Session(engine) as session:
        if date:
            date_begin = datetime.datetime.combine(date, datetime.time.min)
            date_end = date_begin + datetime.timedelta(days=1)
            statement = select(Record).where(Record.date <= date_end).where(Record.date >= date_begin).where(Record.device == str(device))
        else:
            date_begin = datetime.datetime.combine(datetime.datetime.now(), datetime.time.min)
            date_end = date_begin + datetime.timedelta(days=1)
            statement = select(Record).where(Record.date <= date_end).where(Record.date >= date_begin)
            # statement = select(Record)
        results = session.exec(statement)
        resp = []
        for r in results:
            b = str(r.date)[:16]
            r.date = b
            resp.append(r)
    return templates.TemplateResponse(request=request, context={'resp': resp, 'dev': device, 'choosen_date': date}, name="base.html")


@app.get("/example")
async def main_page(request: Request):
    return templates.TemplateResponse(request=request, name="index.html")


@app.get("/chart", response_class=fastapi.responses.HTMLResponse)
async def get_chart(request: Request):
    with Session(engine) as session:
        statement = select(Record)
        results = session.exec(statement)
        resp = []
        for r in results:
            resp.append(r)
        # print(jsonable_encoder(resp))
    return templates.TemplateResponse("chart.html", {"request": request, "json_data": jsonable_encoder(resp)})


@app.get("/chart1")
async def get_e(request: Request):
    return templates.TemplateResponse("ch.html", {"request": request})
