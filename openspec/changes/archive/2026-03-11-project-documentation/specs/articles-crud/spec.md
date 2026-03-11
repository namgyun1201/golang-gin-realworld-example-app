# 아티클 CRUD 및 목록 관리

## ADDED Requirements

### Requirement: 아티클 생성

시스템은 `POST /api/articles` 엔드포인트를 통해 새로운 아티클 생성을 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 유효한 정보로 아티클 생성 성공

- **WHEN** 인증된 사용자가 유효한 `title`(최소 4자), `description`(최대 2048자), `body`(최대 2048자)를 포함한 `article` 객체를 `POST /api/articles`로 전송하면
- **THEN** 시스템은 HTTP 201 (Created) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `article` 객체를 포함해야 하며, `title`, `slug`, `description`, `body`, `createdAt`, `updatedAt`, `author`, `tagList`, `favorited`, `favoritesCount` 필드가 포함되어야 한다(SHALL)
- **THEN** `author` 필드는 작성자의 프로필 정보(`username`, `bio`, `image`, `following`)를 포함해야 한다(SHALL)

#### Scenario: 슬러그 자동 생성

- **WHEN** 아티클이 생성되면
- **THEN** 시스템은 `title`로부터 `slug` 라이브러리(`github.com/gosimple/slug`)를 사용하여 URL 친화적인 슬러그를 자동 생성해야 한다(SHALL)
- **THEN** 슬러그는 데이터베이스에서 고유 인덱스(`uniqueIndex`)로 관리되어야 한다(MUST)

#### Scenario: 태그 목록과 함께 아티클 생성

- **WHEN** `tagList` 배열이 포함된 아티클 생성 요청이 전송되면
- **THEN** 시스템은 기존에 존재하는 태그는 재사용하고 존재하지 않는 태그는 새로 생성해야 한다(SHALL)
- **THEN** 태그와 아티클의 관계는 다대다(many2many) 관계로 `article_tags` 테이블에 저장되어야 한다(SHALL)

#### Scenario: 유효성 검증 실패 시 아티클 생성 실패

- **WHEN** `title`, `description`, `body` 중 하나 이상이 누락되거나 유효성 검증 규칙을 위반하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)

#### Scenario: ArticleUserModel 자동 생성

- **WHEN** 사용자가 처음으로 아티클을 생성하면
- **THEN** 시스템은 `ArticleUserModel`을 `FirstOrCreate`로 자동 생성하여 `UserModel`과 연결해야 한다(SHALL)

### Requirement: 아티클 단건 조회

시스템은 `GET /api/articles/:slug` 엔드포인트를 통해 개별 아티클을 조회할 수 있어야 한다(SHALL).

이 엔드포인트는 인증 없이 접근 가능하다(SHALL).

#### Scenario: 슬러그로 아티클 조회 성공

- **WHEN** 유효한 슬러그로 `GET /api/articles/:slug` 요청이 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `article` 객체를 포함해야 한다(SHALL)
- **THEN** `Author.UserModel`과 `Tags`가 프리로드(Preload)되어 응답에 포함되어야 한다(SHALL)
- **THEN** `createdAt`과 `updatedAt`은 UTC 기준 `2006-01-02T15:04:05.999Z` 형식으로 직렬화되어야 한다(SHALL)
- **THEN** `tagList`는 알파벳 순으로 정렬되어야 한다(SHALL)

#### Scenario: 존재하지 않는 슬러그로 조회 실패

- **WHEN** 존재하지 않는 슬러그로 조회 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `articles` 키로 "Invalid slug" 메시지가 포함되어야 한다(SHALL)

### Requirement: 아티클 수정

시스템은 `PUT /api/articles/:slug` 엔드포인트를 통해 아티클 수정을 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 작성자에 의한 아티클 수정 성공

- **WHEN** 인증된 사용자가 자신이 작성한 아티클의 슬러그로 수정 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 수정된 `article` 객체를 반환해야 한다(SHALL)
- **THEN** 제출되지 않은 필드는 기존 값을 유지해야 한다(SHALL)

#### Scenario: 작성자가 아닌 사용자의 수정 시도

- **WHEN** 인증된 사용자가 다른 사용자가 작성한 아티클을 수정하려고 하면
- **THEN** 시스템은 HTTP 403 (Forbidden) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `article` 키로 "you are not the author" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 존재하지 않는 아티클 수정 시도

- **WHEN** 존재하지 않는 슬러그로 수정 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)

### Requirement: 아티클 삭제

시스템은 `DELETE /api/articles/:slug` 엔드포인트를 통해 아티클 삭제를 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 작성자에 의한 아티클 삭제 성공

- **WHEN** 인증된 사용자가 자신이 작성한 아티클의 슬러그로 삭제 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `{"article": "delete success"}`를 반환해야 한다(SHALL)

#### Scenario: 작성자가 아닌 사용자의 삭제 시도

- **WHEN** 인증된 사용자가 다른 사용자가 작성한 아티클을 삭제하려고 하면
- **THEN** 시스템은 HTTP 403 (Forbidden) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `article` 키로 "you are not the author" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 존재하지 않는 아티클 삭제

- **WHEN** 존재하지 않는 슬러그로 삭제 요청이 전송되면
- **THEN** 시스템은 멱등적으로 동작하여 삭제 작업을 수행하고 성공 응답을 반환해야 한다(SHALL)

### Requirement: 아티클 목록 조회

시스템은 `GET /api/articles` 엔드포인트를 통해 아티클 목록을 조회할 수 있어야 한다(SHALL).

이 엔드포인트는 인증 없이 접근 가능하다(SHALL).

#### Scenario: 전체 아티클 목록 조회

- **WHEN** 필터 없이 `GET /api/articles` 요청이 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `articles` 배열과 `articlesCount` 정수를 포함해야 한다(SHALL)
- **THEN** 기본 limit는 20, 기본 offset은 0으로 적용되어야 한다(SHALL)

#### Scenario: 태그로 아티클 필터링

- **WHEN** `tag` 쿼리 파라미터가 포함된 요청이 전송되면
- **THEN** 시스템은 해당 태그와 연관된 아티클만 반환해야 한다(SHALL)
- **THEN** `TagModel`의 `ArticleModels` 연관을 통해 조회하며, `updated_at desc` 순서로 정렬해야 한다(SHALL)

#### Scenario: 작성자로 아티클 필터링

- **WHEN** `author` 쿼리 파라미터가 포함된 요청이 전송되면
- **THEN** 시스템은 해당 사용자명의 작성자가 작성한 아티클만 반환해야 한다(SHALL)
- **THEN** `UserModel`에서 사용자를 찾고, `ArticleUserModel`의 `ArticleModels` 연관을 통해 조회해야 한다(SHALL)

#### Scenario: 즐겨찾기한 사용자로 아티클 필터링

- **WHEN** `favorited` 쿼리 파라미터가 포함된 요청이 전송되면
- **THEN** 시스템은 해당 사용자명의 사용자가 즐겨찾기한 아티클만 반환해야 한다(SHALL)
- **THEN** `FavoriteModel`에서 해당 사용자의 즐겨찾기를 조회하고 일괄 로드해야 한다(SHALL)

#### Scenario: 페이지네이션 적용

- **WHEN** `limit`과 `offset` 쿼리 파라미터가 포함된 요청이 전송되면
- **THEN** 시스템은 지정된 offset부터 limit 개수만큼의 아티클을 반환해야 한다(SHALL)
- **THEN** `articlesCount`는 필터 조건에 맞는 전체 아티클 수를 반환해야 한다(SHALL)

#### Scenario: 잘못된 limit/offset 값 처리

- **WHEN** `limit`이나 `offset`이 숫자가 아닌 값이면
- **THEN** 시스템은 기본값(limit=20, offset=0)을 적용해야 한다(SHALL)

### Requirement: 아티클 피드

시스템은 `GET /api/articles/feed` 엔드포인트를 통해 팔로우 중인 사용자의 아티클 피드를 제공해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 팔로우 중인 사용자의 아티클 피드 조회

- **WHEN** 인증된 사용자가 `GET /api/articles/feed` 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 팔로우 중인 사용자들이 작성한 아티클의 `articles` 배열과 `articlesCount`를 포함해야 한다(SHALL)
- **THEN** 사용자의 팔로잉 목록에서 `UserModel` ID를 일괄 조회하고, 해당 `ArticleUserModel`들의 아티클을 `updated_at desc` 순서로 반환해야 한다(SHALL)

#### Scenario: 피드 페이지네이션

- **WHEN** `limit`과 `offset` 쿼리 파라미터가 포함된 피드 요청이 전송되면
- **THEN** 시스템은 지정된 페이지네이션을 적용해야 한다(SHALL)
- **THEN** 기본값은 limit=20, offset=0이다(SHALL)

#### Scenario: 팔로잉이 없는 사용자의 피드

- **WHEN** 아무도 팔로우하지 않는 사용자가 피드를 요청하면
- **THEN** 시스템은 빈 `articles` 배열과 `articlesCount: 0`을 반환해야 한다(SHALL)

#### Scenario: 인증되지 않은 사용자의 피드 요청

- **WHEN** 인증 토큰 없이 피드 요청이 전송되면
- **THEN** 시스템은 HTTP 401 (Unauthorized) 상태 코드를 반환해야 한다(SHALL)
