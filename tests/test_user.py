import requests
from uuid import uuid4


base_url = 'http://localhost:8081/api/user'

session_url = f'{base_url}/session'
account_url = f'{base_url}/account'

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


def test_register_ok() -> None:
    creds = {
        'username': str(uuid4()),
        'password': 'bar',
    }

    r = requests.post(account_url, json=creds).json()['data']
    assert r['username'] == creds['username']

    r = requests.get(f'{account_url}/{r["account_id"]}')
    assert r.status_code == 204

    r = requests.post(session_url, json=creds).json()['data']
    assert r['username'] == creds['username']


def test_register_duplicate_username() -> None:
    creds = {
        'username': str(uuid4()),
        'password': 'bar',
    }

    r = requests.post(account_url, json=creds).json()['data']
    assert r['username'] == creds['username']

    r = requests.post(account_url, json=creds).json()
    assert r == {'error': 'username-exists'}


def test_get_exists_account() -> None:
    r = requests.get(f'{account_url}/{exists_author_account_id}')

    assert r.status_code == 204


def test_get_not_exists_account() -> None:
    r = requests.get(f'{account_url}/{str(uuid4())}')

    assert r.status_code == 404


def test_register_for_account_ok() -> None:
    creds = {
        'username': str(uuid4()),
        'password': 'bar',
    }

    r = requests.post(f'{account_url}/{exists_author_account_id}', json=creds).json()['data']
    assert r['username'] == creds['username']

    r = requests.get(f'{account_url}/{r["account_id"]}')
    assert r.status_code == 204

    r = requests.post(session_url, json=creds).json()['data']
    assert r['username'] == creds['username']


def test_register_for_not_exists_account() -> None:
    creds = {
        'username': str(uuid4()),
        'password': 'bar',
    }

    r = requests.post(f'{account_url}/{str(uuid4())}', json=creds).json()
    assert r == {'error': 'account-not-found'}
