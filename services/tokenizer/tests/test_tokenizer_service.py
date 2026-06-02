import unittest
from unittest.mock import patch

from tokenizer.v1 import tokenizer_pb2
from tokenizer_service.server import TokenizerService


class TestTokenizerService(unittest.TestCase):
    @patch("tokenizer_service.tokenize.word_tokenize")
    def test_tokenize(self, mock_word_tokenize):
        mock_word_tokenize.return_value = ["t1", "t2"]
        tokenizer_service = TokenizerService()
        mock_tokens = [
            tokenizer_pb2.Token(token="t1", start=0, end=2),
            tokenizer_pb2.Token(token="t2", start=3, end=5),
        ]
        expected_res = tokenizer_pb2.TokenizeResponse(tokens=mock_tokens)

        req = tokenizer_pb2.TokenizeRequest(text="t1 t2")
        res = tokenizer_service.tokenize(req, None)

        self.assertEqual(res, expected_res)
