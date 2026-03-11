# 사용자 프로필 관리

## Purpose

사용자 프로필 조회, 팔로우/언팔로우 기능의 동작 규격을 정의한다.

## Requirements

### Requirement: 프로필 조회

시스템은 `GET /api/profiles/:username` 엔드포인트를 통해 사용자 프로필을 조회할 수 있어야 한다(SHALL).

#### Scenario: 존재하는 사용자의 프로필 조회 성공

- **WHEN** 존재하는 사용자명으로 `GET /api/profiles/:username` 요청이 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `profile` 객체를 포함해야 하며, `username`, `bio`, `image`, `following` 필드가 포함되어야 한다(SHALL)
- **THEN** `image` 필드는 사용자의 이미지가 설정되지 않은 경우(`nil`) 빈 문자열을 반환해야 한다(SHALL)
- **THEN** `following` 필드는 현재 인증된 사용자가 해당 프로필 사용자를 팔로우하고 있는지 여부를 나타내야 한다(SHALL)

#### Scenario: 인증되지 않은 사용자의 프로필 조회

- **WHEN** 인증 토큰 없이 프로필 조회 요청이 전송되면
- **THEN** 시스템은 `following` 필드를 `false`로 설정하여 프로필을 반환해야 한다(SHALL)

#### Scenario: 존재하지 않는 사용자의 프로필 조회

- **WHEN** 존재하지 않는 사용자명으로 프로필 조회 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `profile` 키로 "Invalid username" 메시지가 포함되어야 한다(SHALL)

### Requirement: 사용자 팔로우

시스템은 `POST /api/profiles/:username/follow` 엔드포인트를 통해 사용자 팔로우 기능을 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 팔로우 성공

- **WHEN** 인증된 사용자가 존재하는 다른 사용자명으로 `POST /api/profiles/:username/follow` 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `profile` 객체를 포함해야 하며, `following` 필드가 `true`로 설정되어야 한다(SHALL)
- **THEN** 팔로우 관계는 `FollowModel` 테이블에 `FollowingID`(대상 사용자)와 `FollowedByID`(현재 사용자)로 저장되어야 한다(SHALL)

#### Scenario: 이미 팔로우한 사용자를 다시 팔로우

- **WHEN** 이미 팔로우 중인 사용자에게 다시 팔로우 요청을 전송하면
- **THEN** 시스템은 `FirstOrCreate`를 사용하여 중복 레코드를 생성하지 않고 정상 응답을 반환해야 한다(SHALL)

#### Scenario: 존재하지 않는 사용자 팔로우 실패

- **WHEN** 존재하지 않는 사용자명으로 팔로우 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `profile` 키로 "Invalid username" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 데이터베이스 오류 시 팔로우 실패

- **WHEN** 팔로우 관계 저장 중 데이터베이스 오류가 발생하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `database` 키로 오류 메시지가 포함되어야 한다(SHALL)

### Requirement: 사용자 언팔로우

시스템은 `DELETE /api/profiles/:username/follow` 엔드포인트를 통해 사용자 언팔로우 기능을 지원해야 한다(SHALL).

이 엔드포인트는 인증이 필수이다(MUST).

#### Scenario: 언팔로우 성공

- **WHEN** 인증된 사용자가 팔로우 중인 사용자명으로 `DELETE /api/profiles/:username/follow` 요청을 전송하면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `profile` 객체를 포함해야 하며, `following` 필드가 `false`로 설정되어야 한다(SHALL)
- **THEN** `FollowModel` 테이블에서 해당 팔로우 관계가 `following_id`와 `followed_by_id` 조건으로 삭제되어야 한다(SHALL)

#### Scenario: 존재하지 않는 사용자 언팔로우 실패

- **WHEN** 존재하지 않는 사용자명으로 언팔로우 요청이 전송되면
- **THEN** 시스템은 HTTP 404 (Not Found) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `profile` 키로 "Invalid username" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 데이터베이스 오류 시 언팔로우 실패

- **WHEN** 언팔로우 관계 삭제 중 데이터베이스 오류가 발생하면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)

### Requirement: 팔로잉 목록 조회

시스템은 사용자가 팔로우하는 사용자 목록을 조회할 수 있어야 한다(SHALL).

#### Scenario: 팔로잉 목록 조회 성공

- **WHEN** 사용자의 팔로잉 목록이 요청되면
- **THEN** 시스템은 `FollowModel` 테이블에서 `FollowedByID`가 현재 사용자인 레코드를 조회해야 한다(SHALL)
- **THEN** `Following` 관계를 Preload하여 팔로잉 중인 `UserModel` 목록을 반환해야 한다(SHALL)
