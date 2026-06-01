import logging
import os
from concurrent import futures

import grpc
from tokenizer.v1 import tokenizer_pb2, tokenizer_pb2_grpc

from tokenizer_service.tokenize import tokenize

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

PORT = os.getenv("TOKENIZER_PORT", "50051")


class TokenizerService(tokenizer_pb2_grpc.TokenizerServiceServicer):
    def tokenize(self, request: tokenizer_pb2.TokenizeRequest, context):
        tokenized = tokenize(request.text)
        tokens = []
        start = 0
        for t in tokenized:
            end = start + len(t)
            token = tokenizer_pb2.Token(token=t, start=start, end=end)
            start = end + 1
            tokens.append(token)

        return tokenizer_pb2.TokenizeResponse(tokens=tokens)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    tokenizer_pb2_grpc.add_TokenizerServiceServicer_to_server(
        TokenizerService(), server
    )

    server.add_insecure_port(f"[::]:{PORT}")
    logger.info(f"Starting TokenizerService on port {PORT}")
    server.start()
    server.wait_for_termination()
