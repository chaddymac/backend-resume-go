import pytest
import unittest
from Visitorcountread import lambda_handler
import json
      
# class TestVisitor(unittest.TestCase):
#     def test_lambda_handler(self):
#         res_api = lambda_handler('a', 'b')
#         self.assertEqual(res_api['statusCode'],200,msg=None)
        
    #need add to pipeline python3 -m unittest test.py


def test_lambda_handler():
    res_api = lambda_handler('a', 'b')
    assert res_api["statusCode"] == 200
    


# def test_lambda_handler(apigw_event, mocker):

#     ret = Visitorcountread.lambda_handler(apigw_event, "")
#     data = json.loads(ret["body"])

#     assert ret["statusCode"] == 200



    # assert "location" in data.dict_keys()



if __name__ == "__main__":
    lambda_handler()