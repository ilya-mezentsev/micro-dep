from typing import Any, Optional
from uuid import uuid4

import requests
import lorem


base_url = 'http://localhost:8080/api/dependencies'

entities_url = f'{base_url}/entities'
entity_url = f'{base_url}/entity'

relation_url = f'{base_url}/relation'

auth_cookie = {
    'X-Auth-Token': 'f93676f8f379c74cefc0d9959d843ac0',
}


def test_get_all_entities() -> None:
    assert len(read_all_entities()) > 0


def test_read_one_entity() -> None:
    all_entities = read_all_entities()
    concrete_entity = read_one_entity(all_entities[0].get('id'))

    assert concrete_entity == all_entities[0]


def test_create_entity() -> None:
    all_entities = read_all_entities()
    author_id = all_entities[0]['author_id']

    created_entity = create_entity({
        'author_id': author_id,
        'name': str(uuid4()),
        'description': lorem.paragraph(),
        'endpoints': [
            {
                'kind': 'foo',
                'address': 'bar',
            },
            {
                'kind': 'foo42',
                'address': 'bar42',
            },
        ],
    })

    concrete_entity = read_one_entity(created_entity['id'])

    assert concrete_entity == created_entity


def test_create_entity_no_endpoints() -> None:
    all_entities = read_all_entities()
    author_id = all_entities[0]['author_id']

    created_entity = create_entity({
        'author_id': author_id,
        'name': str(uuid4()),
        'description': lorem.paragraph(),
    })

    concrete_entity = read_one_entity(created_entity['id'])

    assert concrete_entity == created_entity


def test_create_entity_with_exists_name() -> None:
    all_entities = read_all_entities()
    author_id = all_entities[0]['author_id']
    exists_name = all_entities[0]['name']

    response = create_entity_expect_error({
        'author_id': author_id,
        'name': exists_name,
        'description': lorem.paragraph(),
    })

    assert response['error'] == 'already-exists'


def test_update_entity() -> None:
    all_entities = read_all_entities()
    updated_entity = update_entity({
        **all_entities[0],
        'description': lorem.paragraph(),
    })
    concrete_entity = read_one_entity(all_entities[0]['id'])

    assert concrete_entity == updated_entity


def test_update_entity_remove_endpoint_in_use() -> None:
    all_entities = read_all_entities()
    entity_with_endpoints, related_entity = find_entities_for_relation(all_entities)

    relation = create_relation(related_entity['id'], entity_with_endpoints['endpoints'][0]['id'])

    try:
        response = update_entity_expect_error({
            **entity_with_endpoints,
            'endpoints': entity_with_endpoints['endpoints'][1:]
        })

        assert response['error'] == 'trying-to-remove-endpoint-that-is-in-use'
    finally:
        delete_relation(relation['id'])


def test_delete_entity() -> None:
    all_entities = read_all_entities()

    deleted_id = all_entities[0]['id']
    delete_entity(deleted_id)

    for entity in read_all_entities():
        assert entity['id'] != deleted_id


def test_delete_entity_in_use() -> None:
    all_entities = read_all_entities()
    entity_with_endpoints, related_entity = find_entities_for_relation(all_entities)

    relation = create_relation(related_entity['id'], entity_with_endpoints['endpoints'][0]['id'])

    try:
        response = delete_entity_expect_error(entity_with_endpoints['id'])

        assert response['error'] == 'trying-to-remove-entity-that-is-in-use'
    finally:
        delete_relation(relation['id'])


def read_all_entities() -> list[dict[str, Any]]:
    entities = requests.get(entities_url, cookies=auth_cookie).json()['data']
    for i, entity in enumerate(entities):
        entities[i] = sort_endpoints(entity)

    return sorted(entities, key=lambda e: e['id'])


def read_one_entity(id_: str) -> dict[str, Any]:
    return sort_endpoints(requests.get(f'{entity_url}/{id_}', cookies=auth_cookie).json()['data'])


def create_entity(entity: dict[str, Any]) -> dict[str, Any]:
    response = requests.post(entity_url, json=entity, cookies=auth_cookie).json()

    assert 'data' in response, f'no data in response: {response}'

    return sort_endpoints(response['data'])


def create_entity_expect_error(entity: dict[str, Any]) -> dict[str, Any]:
    return requests.post(entity_url, json=entity, cookies=auth_cookie).json()


def update_entity(entity: dict[str, Any]) -> dict[str, Any]:
    response = requests.put(entity_url, json=entity, cookies=auth_cookie).json()

    assert 'data' in response, f'no data in response: {response}'

    return sort_endpoints(response['data'])


def update_entity_expect_error(entity: dict[str, Any]) -> dict[str, Any]:
    return requests.put(entity_url, json=entity, cookies=auth_cookie).json()


def delete_entity(id_: str) -> None:
    requests.delete(f'{entity_url}/{id_}', cookies=auth_cookie)


def delete_entity_expect_error(id_: str) -> dict[str, Any]:
    return requests.delete(f'{entity_url}/{id_}', cookies=auth_cookie).json()


def create_relation(from_entity_id: str, to_endpoint_id: str) -> dict[str, Any]:
    response = requests.post(
        relation_url,
        json={'from_entity_id': from_entity_id, 'to_endpoint_id': to_endpoint_id},
        cookies=auth_cookie,
    ).json()

    assert 'data' in response, f'no data in response: {response}'

    return response['data']


def delete_relation(id_: str) -> None:
    requests.delete(f'{relation_url}/{id_}', cookies=auth_cookie)


def sort_endpoints(entity: dict[str, Any]) -> dict[str, Any]:
    if not isinstance(entity['endpoints'], list):
        entity['endpoints'] = []

    entity['endpoints'] = sorted(entity['endpoints'], key=lambda e: e['id'])

    return entity


def find_entities_for_relation(all_entities: list[dict[str, Any]]) -> tuple[dict[str, Any], dict[str, Any]]:
    # find entity with endpoints for relation creation
    entity_with_endpoints: Optional[dict[str, Any]] = None
    for entity in all_entities:
        if entity['endpoints']:
            entity_with_endpoints = entity
            break

    assert entity_with_endpoints is not None, 'unable to find entity with endpoints'

    # find another entity for create relation to one of entity_with_endpoints's endpoint
    related_entity: Optional[dict[str, Any]] = None
    for entity in all_entities:
        if entity['id'] != entity_with_endpoints['id']:
            related_entity = entity
            break

    assert related_entity is not None, 'unable to find related entity'

    return entity_with_endpoints, related_entity
