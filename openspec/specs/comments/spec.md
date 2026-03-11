# 댓글 관리

## Purpose

아티클 댓글 생성/조회/삭제 기능의 동작 규격을 정의한다.

## Requirements

### Requirement: 댓글 생성

시스템은 `POST /api/articles/:slug/comments` 엔드포인트를 통해 아티클에 댓글을 추가할 수 있어야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 유효한 댓글 생성 성공

- **WHEN** 인증된 사용자가 존재하는 아티클의 슬러그와 유효한 `body`(최대 2048자)를 포함한 `comment` 객체를 `POST /api/articles/:slug/comments`로 전송하면
- **THEN** 시스템은 HTTP 201 (Created) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `comment` 객체를 포함해야 하며, `id`, `body`, `createdAt`, `updatedAt`, `author` 필드가 포함되어야 한다(SHALL)
- **THEN** `author` 필드는 작성자의 프로필 정보(`username`, `bio`, `image`, `following`)를 포함해야 한다(SHALL)
- **THEN** `createdAt`과 `updatedAt`은 UTC 기준 `2006-01-02T15:04:05.999Z` 형식으로 직렬화되어야 한다(SHALL)

#### Scenario: 댓글의 아티클 연관 설정

- **WHEN** 댓글이 생성되면
- **THEN** 시스템은 `CommentModel`의 `Article` 필드에 해당 아티클 모델을 설정해야 한다(SHALL)
- **THEN** 댓글의 `Author`는 현재 인증된 사용자의 `ArticleUserModel`로 설정되어야 한다(SHALL)

#### Scenario: 존재하지 않는 아티클에 댓글 생성 실패

- **WHEN** 존재하지 않는 슬러그로 댓글 생성 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `comment` 키로 "Invalid slug" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 유효성 검증 실패 시 댓글 생성 실패

- **WHEN** `body` 필드가 누락되거나 2048자를 초과하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문에 유효성 검증 오류 정보가 포함되어야 한다(SHALL)

#### Scenario: 데이터베이스 오류 시 댓글 생성 실패

- **WHEN** 댓글 저장 중 데이터베이스 오류가 발생하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `database` 키로 오류 메시지가 포함되어야 한다(SHALL)

### Requirement: 댓글 목록 조회

시스템은 `GET /api/articles/:slug/comments` 엔드포인트를 통해 아티클의 댓글 목록을 조회할 수 있어야 한다(SHALL).

이 엔드포인트는 인증 없이 접근 가능하다(SHALL).

#### Scenario: 아티클의 댓글 목록 조회 성공

- **WHEN** 존재하는 아티클의 슬러그로 `GET /api/articles/:slug/comments` 요청이 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `comments` 배열을 포함해야 하며, 각 댓글은 `id`, `body`, `createdAt`, `updatedAt`, `author` 필드를 가져야 한다(SHALL)
- **THEN** 각 댓글의 `Author.UserModel`이 프리로드(Preload)되어야 한다(SHALL)

#### Scenario: 댓글이 없는 아티클의 목록 조회

- **WHEN** 댓글이 없는 아티클의 댓글 목록을 조회하면
- **THEN** 시스템은 빈 `comments` 배열을 반환해야 한다(SHALL)

#### Scenario: 존재하지 않는 아티클의 댓글 목록 조회

- **WHEN** 존재하지 않는 슬러그로 댓글 목록 조회 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `comments` 키로 "Invalid slug" 메시지가 포함되어야 한다(SHALL)

### Requirement: 댓글 삭제

시스템은 `DELETE /api/articles/:slug/comments/:id` 엔드포인트를 통해 댓글 삭제를 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 작성자에 의한 댓글 삭제 성공

- **WHEN** 인증된 사용자가 자신이 작성한 댓글의 ID로 `DELETE /api/articles/:slug/comments/:id` 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `{"comment": "delete success"}`를 반환해야 한다(SHALL)

#### Scenario: 작성자가 아닌 사용자의 댓글 삭제 시도

- **WHEN** 인증된 사용자가 다른 사용자가 작성한 댓글을 삭제하려고 하면
- **THEN** 시스템은 HTTP 403 (Forbidden) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `comment` 키로 "you are not the author" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 존재하지 않는 댓글 삭제

- **WHEN** 존재하지 않는 댓글 ID로 삭제 요청이 전송되면
- **THEN** 시스템은 멱등적으로 동작하여 삭제 작업을 수행하고 성공 응답을 반환해야 한다(SHALL)

#### Scenario: 잘못된 댓글 ID 형식

- **WHEN** 숫자가 아닌 댓글 ID로 삭제 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `comment` 키로 "Invalid id" 메시지가 포함되어야 한다(SHALL)
