# 즐겨찾기 및 태그 관리

## Purpose

아티클 즐겨찾기 추가/해제, 태그 목록 조회 기능의 동작 규격을 정의한다.

## Requirements

### Requirement: 즐겨찾기 추가

시스템은 `POST /api/articles/:slug/favorite` 엔드포인트를 통해 아티클 즐겨찾기 추가를 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 아티클 즐겨찾기 추가 성공

- **WHEN** 인증된 사용자가 존재하는 아티클의 슬러그로 `POST /api/articles/:slug/favorite` 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `article` 객체를 반환해야 하며, `favorited` 필드가 `true`로, `favoritesCount`가 증가된 값으로 설정되어야 한다(SHALL)
- **THEN** 즐겨찾기 관계는 `FavoriteModel` 테이블에 `FavoriteID`(아티클 ID)와 `FavoriteByID`(ArticleUserModel ID)로 저장되어야 한다(SHALL)

#### Scenario: 이미 즐겨찾기한 아티클에 다시 추가

- **WHEN** 이미 즐겨찾기한 아티클에 다시 즐겨찾기 요청을 전송하면
- **THEN** 시스템은 `FirstOrCreate`를 사용하여 중복 레코드를 생성하지 않고 정상 응답을 반환해야 한다(SHALL)

#### Scenario: 존재하지 않는 아티클 즐겨찾기 실패

- **WHEN** 존재하지 않는 슬러그로 즐겨찾기 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `articles` 키로 "Invalid slug" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 데이터베이스 오류 시 즐겨찾기 추가 실패

- **WHEN** 즐겨찾기 저장 중 데이터베이스 오류가 발생하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)

### Requirement: 즐겨찾기 해제

시스템은 `DELETE /api/articles/:slug/favorite` 엔드포인트를 통해 아티클 즐겨찾기 해제를 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 아티클 즐겨찾기 해제 성공

- **WHEN** 인증된 사용자가 즐겨찾기한 아티클의 슬러그로 `DELETE /api/articles/:slug/favorite` 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `article` 객체를 반환해야 하며, `favorited` 필드가 `false`로, `favoritesCount`가 감소된 값으로 설정되어야 한다(SHALL)
- **THEN** `FavoriteModel` 테이블에서 해당 즐겨찾기 관계가 `favorite_id`와 `favorite_by_id` 조건으로 삭제되어야 한다(SHALL)

#### Scenario: 존재하지 않는 아티클 즐겨찾기 해제 실패

- **WHEN** 존재하지 않는 슬러그로 즐겨찾기 해제 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)

#### Scenario: 데이터베이스 오류 시 즐겨찾기 해제 실패

- **WHEN** 즐겨찾기 삭제 중 데이터베이스 오류가 발생하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)

### Requirement: 아티클의 즐겨찾기 수

시스템은 각 아티클의 즐겨찾기 수(`favoritesCount`)를 정확하게 계산해야 한다(SHALL).

#### Scenario: 단건 아티클의 즐겨찾기 수 조회

- **WHEN** 단건 아티클이 조회되면
- **THEN** 시스템은 `FavoriteModel` 테이블에서 해당 아티클의 `FavoriteID`로 `COUNT` 쿼리를 실행하여 즐겨찾기 수를 계산해야 한다(SHALL)
- **THEN** `favoritesCount` 필드에 `uint` 타입으로 반환해야 한다(SHALL)

#### Scenario: 아티클 목록의 즐겨찾기 수 일괄 조회

- **WHEN** 아티클 목록이 조회되면
- **THEN** 시스템은 `BatchGetFavoriteCounts`를 사용하여 모든 아티클의 즐겨찾기 수를 단일 쿼리(`GROUP BY favorite_id`)로 일괄 조회해야 한다(SHALL)
- **THEN** N+1 쿼리 문제를 방지하기 위해 `BatchGetFavoriteStatus`를 사용하여 현재 사용자의 즐겨찾기 상태도 일괄 조회해야 한다(SHALL)

#### Scenario: 현재 사용자의 즐겨찾기 여부

- **WHEN** 인증된 사용자가 아티클을 조회하면
- **THEN** `favorited` 필드는 현재 사용자의 `ArticleUserModel` ID와 아티클 ID로 `FavoriteModel` 테이블을 조회하여 결정해야 한다(SHALL)

#### Scenario: 인증되지 않은 사용자의 즐겨찾기 여부

- **WHEN** 인증되지 않은 사용자가 아티클을 조회하면
- **THEN** `favorited` 필드는 `false`로 설정되어야 한다(SHALL)

### Requirement: 태그 목록 조회

시스템은 `GET /api/tags` 엔드포인트를 통해 전체 태그 목록을 조회할 수 있어야 한다(SHALL).

이 엔드포인트는 인증 없이 접근 가능하다(SHALL).

#### Scenario: 전체 태그 목록 조회 성공

- **WHEN** `GET /api/tags` 요청이 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `tags` 배열을 포함해야 하며, 각 요소는 태그 문자열이어야 한다(SHALL)
- **THEN** 데이터베이스의 모든 `TagModel` 레코드에서 `Tag` 필드를 추출하여 반환해야 한다(SHALL)

#### Scenario: 태그가 없는 경우

- **WHEN** 데이터베이스에 태그가 없을 때 태그 목록을 조회하면
- **THEN** 시스템은 빈 `tags` 배열을 반환해야 한다(SHALL)

#### Scenario: 태그 목록 조회 실패

- **WHEN** 태그 목록 조회 중 데이터베이스 오류가 발생하면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `articles` 키로 "Invalid param" 메시지가 포함되어야 한다(SHALL)

### Requirement: 태그의 고유성

시스템은 태그의 고유성을 보장해야 한다(MUST).

#### Scenario: 태그 고유 인덱스

- **WHEN** 새로운 태그가 생성되면
- **THEN** `TagModel`의 `Tag` 필드는 데이터베이스에서 고유 인덱스(`uniqueIndex`)로 관리되어야 한다(MUST)

#### Scenario: 동시 태그 생성 시 중복 방지

- **WHEN** 동시에 동일한 태그를 생성하려고 하면
- **THEN** 시스템은 생성 실패 시 기존 태그를 조회하여 사용하는 방식으로 경쟁 조건을 처리해야 한다(SHALL)
