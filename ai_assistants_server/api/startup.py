from .routers import azure_openai
from fastapi import FastAPI

def get_app(lifespan) -> FastAPI:
    app = FastAPI(lifespan=lifespan)

    # Include the Azure OpenAI router
    app.include_router(azure_openai.router, prefix="/api/v1")

    return app
