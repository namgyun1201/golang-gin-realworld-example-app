# 사용자 인증 및 계정 관리

## ADDED Requirements

### Requirement: 사용자 등록

시스템은 `POST /api/users` 엔드포인트를 통해 새로운 사용자 등록을 지원해야 한다(SHALL).

요청 본문은 `user` 객체 내에 `username`, `email`, `password` 필드를 포함해야 한다(MUST).

#### Scenario: 유효한 정보로 사용자 등록 성공

- **WHEN** 유효한 `username`(최소 4자, 최대 255자), `email`(이메일 형식), `password`(최소 8자, 최대 255자)를 포함한 요청이 `POST /api/users`로 전송되면
- **THEN** 시스템은 HTTP 201 (Created) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `user` 객체를 포함해야 하며, `username`, `email`, `bio`, `image`, `token` 필드가 포함되어야 한다(SHALL)
- **THEN** `token` 필드에는 유효한 JWT 토큰이 포함되어야 한다(SHALL)
- **THEN** 비밀번호는 bcrypt로 해싱되어 저장되어야 한다(MUST)

#### Scenario: 필수 필드 누락 시 등록 실패

- **WHEN** `username`, `email`, `password` 중 하나 이상이 누락된 요청이 전송되면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `errors` 객체를 포함해야 하며, 각 필드별 유효성 검증 오류 메시지가 포함되어야 한다(SHALL)

#### Scenario: 유효성 검증 규칙 위반 시 등록 실패

- **WHEN** `username`이 4자 미만이거나, `email`이 이메일 형식이 아니거나, `password`가 8자 미만인 요청이 전송되면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 해당 필드의 유효성 검증 태그와 파라미터 정보가 포함되어야 한다(SHALL)

#### Scenario: 중복 이메일로 등록 실패

- **WHEN** 이미 등록된 이메일 주소로 등록 요청이 전송되면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 데이터베이스 오류 정보가 포함되어야 한다(SHALL)

### Requirement: 사용자 로그인

시스템은 `POST /api/users/login` 엔드포인트를 통해 사용자 로그인을 지원해야 한다(SHALL).

요청 본문은 `user` 객체 내에 `email`, `password` 필드를 포함해야 한다(MUST).

#### Scenario: 유효한 자격 증명으로 로그인 성공

- **WHEN** 등록된 이메일과 올바른 비밀번호를 포함한 요청이 `POST /api/users/login`으로 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `user` 객체를 포함해야 하며, `username`, `email`, `bio`, `image`, `token` 필드가 포함되어야 한다(SHALL)
- **THEN** 컨텍스트에 사용자 모델과 사용자 ID가 설정되어야 한다(SHALL)

#### Scenario: 등록되지 않은 이메일로 로그인 실패

- **WHEN** 등록되지 않은 이메일로 로그인 요청이 전송되면
- **THEN** 시스템은 HTTP 401 (Unauthorized) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `login` 키로 "Not Registered email or invalid password" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 잘못된 비밀번호로 로그인 실패

- **WHEN** 올바른 이메일이지만 잘못된 비밀번호로 로그인 요청이 전송되면
- **THEN** 시스템은 HTTP 401 (Unauthorized) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문의 `errors` 객체에 `login` 키로 "Not Registered email or invalid password" 메시지가 포함되어야 한다(SHALL)

#### Scenario: 로그인 유효성 검증 실패

- **WHEN** `email`이 이메일 형식이 아니거나 `password`가 8자 미만인 요청이 전송되면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)

### Requirement: 현재 사용자 조회

시스템은 `GET /api/user` 엔드포인트를 통해 인증된 현재 사용자 정보를 반환해야 한다(SHALL).

#### Scenario: 인증된 사용자 정보 조회 성공

- **WHEN** 유효한 JWT 토큰을 포함한 요청이 `GET /api/user`로 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 `user` 객체를 포함해야 하며, `username`, `email`, `bio`, `image`, `token` 필드가 포함되어야 한다(SHALL)
- **THEN** `image` 필드는 설정되지 않은 경우 빈 문자열을 반환해야 한다(SHALL)

#### Scenario: 인증 토큰 없이 조회 실패

- **WHEN** JWT 토큰 없이 요청이 `GET /api/user`로 전송되면
- **THEN** 시스템은 HTTP 401 (Unauthorized) 상태 코드를 반환해야 한다(SHALL)

### Requirement: 사용자 정보 수정

시스템은 `PUT /api/user` 엔드포인트를 통해 인증된 사용자의 정보 수정을 지원해야 한다(SHALL).

#### Scenario: 사용자 정보 수정 성공

- **WHEN** 유효한 JWT 토큰과 수정할 필드(`username`, `email`, `password`, `bio`, `image`)를 포함한 요청이 `PUT /api/user`로 전송되면
- **THEN** 시스템은 HTTP 200 (OK) 상태 코드를 반환해야 한다(SHALL)
- **THEN** 응답 본문은 수정된 `user` 객체를 반환해야 한다(SHALL)
- **THEN** 변경되지 않은 필드는 기존 값을 유지해야 한다(SHALL)

#### Scenario: 비밀번호 변경 시 해싱 적용

- **WHEN** 새로운 비밀번호가 포함된 수정 요청이 전송되면
- **THEN** 비밀번호가 기존 더미 값(RandomPassword)과 다른 경우에만 bcrypt로 새로 해싱되어야 한다(SHALL)

#### Scenario: 이미지 URL 수정

- **WHEN** `image` 필드에 유효한 URL이 포함된 수정 요청이 전송되면
- **THEN** 시스템은 해당 URL을 사용자의 이미지로 저장해야 한다(SHALL)
- **THEN** `image` 필드는 `url` 형식 유효성 검증을 통과해야 한다(MUST)

#### Scenario: 수정 유효성 검증 실패

- **WHEN** 유효성 검증 규칙을 위반하는 수정 요청이 전송되면
- **THEN** 시스템은 HTTP 422 (Unprocessable Entity) 상태 코드를 반환해야 한다(SHALL)

### Requirement: JWT 토큰 인증 흐름

시스템은 JWT 토큰 기반 인증을 지원해야 한다(SHALL).

#### Scenario: Authorization 헤더를 통한 토큰 전달

- **WHEN** `Authorization: Token <jwt_token>` 형식의 헤더가 포함된 요청이 전송되면
- **THEN** 시스템은 "Token " 접두사(대소문자 무시) 이후의 문자열을 JWT 토큰으로 추출해야 한다(SHALL)

#### Scenario: 쿼리 파라미터를 통한 토큰 전달

- **WHEN** `access_token` 쿼리 파라미터가 포함된 요청이 전송되면
- **THEN** 시스템은 해당 값을 JWT 토큰으로 사용해야 한다(SHALL)

#### Scenario: 유효한 JWT 토큰 검증

- **WHEN** HMAC-SHA256으로 서명된 유효한 JWT 토큰이 전달되면
- **THEN** 시스템은 토큰에서 `id` 클레임을 추출하여 사용자를 식별해야 한다(SHALL)
- **THEN** 컨텍스트에 `my_user_id`와 `my_user_model`을 설정해야 한다(SHALL)

#### Scenario: 유효하지 않은 JWT 토큰 (필수 인증 경로)

- **WHEN** 유효하지 않거나 만료된 JWT 토큰이 인증 필수 경로(`auto401=true`)에 전달되면
- **THEN** 시스템은 HTTP 401 (Unauthorized) 상태 코드를 반환해야 한다(SHALL)

#### Scenario: 토큰 없는 요청 (선택적 인증 경로)

- **WHEN** JWT 토큰 없이 선택적 인증 경로(`auto401=false`)에 요청이 전달되면
- **THEN** 시스템은 `my_user_id`를 0으로, `my_user_model`을 빈 모델로 설정하고 요청을 계속 처리해야 한다(SHALL)

### Requirement: 에러 응답 포맷

모든 에러 응답은 일관된 형식을 따라야 한다(MUST).

#### Scenario: 유효성 검증 에러 응답

- **WHEN** go-playground/validator 유효성 검증이 실패하면
- **THEN** 응답 본문은 `errors` 객체를 포함해야 하며, 각 필드명을 키로 하고 `{tag: param}` 또는 `{key: tag}` 형식의 문자열을 값으로 가져야 한다(SHALL)

#### Scenario: 일반 에러 응답

- **WHEN** 데이터베이스 오류 등 일반적인 에러가 발생하면
- **THEN** 응답 본문은 `errors` 객체를 포함해야 하며, 지정된 키(예: "database", "login")에 에러 메시지 문자열을 값으로 가져야 한다(SHALL)
