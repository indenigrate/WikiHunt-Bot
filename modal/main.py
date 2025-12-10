import modal

# 1. Define the Container Image
def download_model():
    from sentence_transformers import SentenceTransformer
    # Downloads model to the image cache
    SentenceTransformer('BAAI/bge-base-en-v1.5')

image = (
    modal.Image.debian_slim()
    .pip_install("fastapi", "uvicorn", "sentence-transformers", "torch")
    .run_function(download_model)
)

app = modal.App("bge-similarity-api", image=image)

# 2. Define the Serverless Function
@app.function(
    image=image,
    gpu="any",
    scaledown_window=300,  # UPDATED: Replaces container_idle_timeout
)
@modal.asgi_app()
def fastapi_app():
    from fastapi import FastAPI, HTTPException
    from fastapi.middleware.cors import CORSMiddleware
    from pydantic import BaseModel
    from sentence_transformers import SentenceTransformer, util
    from typing import List
    import torch

    web_app = FastAPI()

    web_app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    device = 'cuda' if torch.cuda.is_available() else 'cpu'
    print(f"Loading BAAI/bge-base-en-v1.5 on {device}...")
    model = SentenceTransformer('BAAI/bge-base-en-v1.5', device=device)

    class BulkSimilarityRequest(BaseModel):
        target: str
        inputs: List[str]

    @web_app.get("/healthz")
    def health_check():
        return {"message": "OK"}

    @web_app.get("/")
    def root():
        return {"message": "SBERT Semantic Similarity API (BGE-Base)"}

    @web_app.post("/similarity")
    def compute_bulk_similarity(request: BulkSimilarityRequest):
        try:
            all_texts = [request.target] + request.inputs
            # normalize_embeddings=True is recommended for BGE models
            embeddings = model.encode(all_texts, convert_to_tensor=True, normalize_embeddings=True)
            
            target_embedding = embeddings[0]
            input_embeddings = embeddings[1:]

            similarities = util.cos_sim(target_embedding, input_embeddings).squeeze().tolist()
            if not isinstance(similarities, list):
                similarities = [similarities]
            return {"similarities": similarities}
        except Exception as e:
            raise HTTPException(status_code=500, detail=str(e))

    return web_app