import base64
import re
import requests
from flask import Flask, request, Response

app = Flask(__name__)

@app.route('/replace', methods=['GET'])
def replace_url_content():
    original_url = request.args.get('url')
    old_value = request.args.get('old_value', 'icook.hk')
    new_value = request.args.get('new_value', 'cfip.gay')

    if not original_url:
        return Response("缺少参数 'url'", status=400)

    try:
        response = requests.get(original_url)
        response.raise_for_status()
        encoded_data = response.text
        decoded_data = base64.urlsafe_b64decode(encoded_data).decode('utf-8')
    except Exception as e:
        return Response(f"获取或解码错误: {str(e)}", status=400)

    vless_pattern = re.compile(r'(vless://[A-Za-z0-9-_]+)@')
    trojan_pattern = re.compile(r'(trojan://[A-Za-z0-9-_]+)@')
    replaced_data = decoded_data
    for match in vless_pattern.finditer(decoded_data):
        vless_id = match.group(1)
        replaced_data = replaced_data.replace(f"{vless_id}@{old_value}", f"{vless_id}@{new_value}")
    for match in trojan_pattern.finditer(decoded_data):
        trojan_id = match.group(1)
        replaced_data = replaced_data.replace(f"{trojan_id}@{old_value}", f"{trojan_id}@{new_value}")

    try:
        encoded_replaced_data = base64.b64encode(replaced_data.encode('utf-8')).decode('utf-8')
    except Exception as e:
        return Response(f"编码错误: {str(e)}", status=400)

    return Response(encoded_replaced_data, content_type='text/plain')

if __name__ == '__main__':
    app.run(debug=True)
