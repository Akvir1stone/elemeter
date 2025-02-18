from fastapi import FastAPI
from models import Record


app = FastAPI()


@app.post("/")
async def new_record(record: Record):
    return None


@app.get("/")
async def main_page():
    return None
