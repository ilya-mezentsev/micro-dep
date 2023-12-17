import requests


base_url = 'http://localhost:8081/api/user'

session_url = f'{base_url}/session'

auth_cookie = {
    'X-Auth-Token': 'f93676f8f379c74cefc0d9959d843ac0',
}

exists_author_id = '4a7b8037-2cba-4667-bf1b-a2d227c7b80b'
exists_author_username = 'foo'
exists_author_registered_at = 1699191331
exists_author_account_id = 'c5d6f777-8195-4908-af7d-6e3c1cd2460e'


def test_get_session() -> None:
    r = requests.get(session_url, cookies=auth_cookie).json()['data']

    assert r == {
        'id': exists_author_id,
        'account_id': exists_author_account_id,
        'username': exists_author_username,
        'registered_at': 1699191331,
    }


def test_get_session_wrong_token() -> None:
    r = requests.get(session_url, cookies={
        'X-Auth-Token': 'foo-bar',
    }).json()

    assert r == {'error': 'account-not-found'}


def test_delete_session() -> None:
    r = requests.delete(session_url)

    assert r.status_code == 204
    assert r.headers['Set-Cookie'] == 'X-Auth-Token=; Path=/; Domain=localhost; HttpOnly'
