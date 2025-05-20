from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer, util

app = FastAPI() # starts api 
# model = SentenceTransformer('all-mpnet-base-v2')  # loads sbert model
model = SentenceTransformer('all-MiniLM-L6-v2')  # loads sbert model

class TextPair(BaseModel):
    input1: str
    input2: str

@app.get("/")
def root():
    return {"message": "SBERT Semantic Similarity API"}

@app.post("/similarity")
def compute_similarity(pair: TextPair):
    try:
        embeddings = model.encode([pair.input1, pair.input2])
        similarity = util.cos_sim(embeddings[0], embeddings[1]).item()
        return {"similarity": similarity}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
