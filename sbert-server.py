from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer, util
from typing import List
import torch

app = FastAPI()

# model = SentenceTransformer('all-mpnet-base-v2')
# model = SentenceTransformer('all-MiniLM-L6-v2')
model = SentenceTransformer('all-mpnet-base-v2', device='cuda' if torch.cuda.is_available() else 'cpu')


class BulkSimilarityRequest(BaseModel):
    target: str
    inputs: List[str]

@app.get("/")
def root():
    return {"message": "SBERT Semantic Similarity API"}

@app.post("/similarity")
def compute_bulk_similarity(request: BulkSimilarityRequest):
    try:
        # Combine target and all input strings for batch encoding
        all_texts = [request.target] + request.inputs
        embeddings = model.encode(all_texts, convert_to_tensor=True)
        
        target_embedding = embeddings[0]
        input_embeddings = embeddings[1:]

        # Compute cosine similarities
        similarities = util.cos_sim(target_embedding, input_embeddings).squeeze().tolist()
        return {"similarities": similarities}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
