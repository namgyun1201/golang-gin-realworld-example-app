# 공통 인프라스트럭처

## ADDED Requirements

### Requirement: 데이터베이스 초기화 및 AutoMigrate

시스템은 SQLite 데이터베이스를 사용하여 데이터를 저장하고 스키마를 자동으로 마이그레이션해야 한다(SHALL).

#### Scenario: 프로덕션 데이터베이스 초기화

- **WHEN** `common.Init()` 함수가 호출되면
- **THEN** 시스템은 `DB_PATH` 환경 변수에 지정된 경로 또는 기본값 `./data/gorm.db`에 SQLite 데이터베이스를 생성해야 한다(SHALL)
- **THEN** 데이터베이스 파일의 디렉토리가 존재하지 않으면 `0750` 권한으로 자동 생성해야 한다(SHALL)
- **THEN** 최대 유휴 연결 수를 10으로 설정해야 한다(SHALL)
- **THEN** 글로벌 변수 `DB`에 데이터베이스 연결을 저장해야 한다(SHALL)

#### Scenario: 테스트 데이터베이스 초기화

- **WHEN** `common.TestDBInit()` 함수가 호출되면
- **THEN** 시스템은 `TEST_DB_PATH` 환경 변수에 지정된 경로 또는 기본값 `./data/gorm_test.db`에 테스트 데이터베이스를 생성해야 한다(SHALL)
- **THEN** GORM 로거를 `Info` 레벨로 설정해야 한다(SHALL)
- **THEN** 최대 유휴 연결 수를 3으로 설정해야 한다(SHALL)

#### Scenario: 테스트 데이터베이스 정리

- **WHEN** `common.TestDBFree()` 함수가 호출되면
- **THEN** 시스템은 데이터베이스 연결을 닫고 테스트 데이터베이스 파일을 삭제해야 한다(SHALL)

#### Scenario: 사용자 모델 AutoMigrate

- **WHEN** `users.AutoMigrate()` 함수가 호출되면
- **THEN** 시스템은 `UserModel`과 `FollowModel` 테이블의 스키마를 자동 마이그레이션해야 한다(SHALL)

#### Scenario: 데이터베이스 연결 조회

- **WHEN** `common.GetDB()` 함수가 호출되면
- **THEN** 시스템은 글로벌 변수 `DB`에 저장된 데이터베이스 연결을 반환해야 한다(SHALL)

### Requirement: JWT 토큰 생성 및 검증

시스템은 JWT 토큰을 생성하고 검증하는 유틸리티를 제공해야 한다(SHALL).

#### Scenario: JWT 토큰 생성

- **WHEN** `common.GenToken(id)` 함수가 사용자 ID와 함께 호출되면
- **THEN** 시스템은 HMAC-SHA256(`jwt.SigningMethodHS256`) 서명 방식으로 JWT 토큰을 생성해야 한다(SHALL)
- **THEN** 토큰의 클레임에 `id`(사용자 ID)와 `exp`(현재 시간 + 24시간) 필드를 포함해야 한다(SHALL)
- **THEN** 토큰 서명에 `common.JWTSecret` 상수를 사용해야 한다(MUST)

#### Scenario: JWT 토큰 서명 실패

- **WHEN** 토큰 서명 과정에서 오류가 발생하면
- **THEN** 시스템은 오류 메시지를 출력하고 빈 문자열을 반환해야 한다(SHALL)

#### Scenario: JWT 토큰 검증 (미들웨어)

- **WHEN** `AuthMiddleware`에서 JWT 토큰을 검증하면
- **THEN** 시스템은 서명 방식이 `HMAC`인지 확인해야 한다(SHALL)
- **THEN** 서명 방식이 올바르지 않으면 `jwt.ErrSignatureInvalid` 오류를 반환해야 한다(SHALL)
- **THEN** 유효한 토큰의 `id` 클레임에서 사용자 ID를 추출하여 `float64`에서 `uint`로 변환해야 한다(SHALL)

### Requirement: 입력 유효성 검증

시스템은 `go-playground/validator`를 사용하여 요청 본문의 유효성을 검증해야 한다(SHALL).

#### Scenario: 바인딩 처리

- **WHEN** `common.Bind(c, obj)` 함수가 호출되면
- **THEN** 시스템은 요청의 HTTP 메서드와 Content-Type에 따라 적절한 바인딩 방식(`binding.Default`)을 선택해야 한다(SHALL)
- **THEN** `c.ShouldBindWith`를 사용하여 바인딩을 수행해야 한다(SHALL) (자동 400 응답을 방지하기 위해 `c.MustBindWith` 대신 사용)

#### Scenario: 사용자 등록 유효성 검증 규칙

- **WHEN** 사용자 등록 요청이 바인딩되면
- **THEN** 시스템은 다음 규칙을 적용해야 한다(SHALL):
  - `username`: 필수, 최소 4자, 최대 255자
  - `email`: 필수, 이메일 형식
  - `password`: 필수, 최소 8자, 최대 255자
  - `bio`: 선택, 최대 1024자
  - `image`: 선택, URL 형식 (`omitempty,url`)

#### Scenario: 로그인 유효성 검증 규칙

- **WHEN** 로그인 요청이 바인딩되면
- **THEN** 시스템은 다음 규칙을 적용해야 한다(SHALL):
  - `email`: 필수, 이메일 형식
  - `password`: 필수, 최소 8자, 최대 255자

#### Scenario: 아티클 유효성 검증 규칙

- **WHEN** 아티클 생성/수정 요청이 바인딩되면
- **THEN** 시스템은 다음 규칙을 적용해야 한다(SHALL):
  - `title`: 필수, 최소 4자
  - `description`: 필수, 최대 2048자
  - `body`: 필수, 최대 2048자
  - `tagList`: 선택, 문자열 배열

#### Scenario: 댓글 유효성 검증 규칙

- **WHEN** 댓글 생성 요청이 바인딩되면
- **THEN** 시스템은 다음 규칙을 적용해야 한다(SHALL):
  - `body`: 필수, 최대 2048자

### Requirement: 에러 응답 포맷

시스템은 모든 에러에 대해 일관된 응답 형식을 사용해야 한다(MUST).

#### Scenario: CommonError 구조체

- **WHEN** 에러 응답이 생성되면
- **THEN** 응답 본문은 `errors` 키를 가진 JSON 객체를 포함해야 한다(SHALL)
- **THEN** `errors`의 값은 키-값 쌍의 맵(`map[string]interface{}`)이어야 한다(SHALL)

#### Scenario: 유효성 검증 에러 (NewValidatorError)

- **WHEN** `go-playground/validator`의 `ValidationErrors`가 발생하면
- **THEN** 시스템은 각 유효성 검증 오류를 순회하여 `errors` 맵에 추가해야 한다(SHALL)
- **THEN** 필드명을 키로, 파라미터가 있으면 `{tag: param}` 형식, 없으면 `{key: tag}` 형식의 문자열을 값으로 설정해야 한다(SHALL)

#### Scenario: 일반 에러 (NewError)

- **WHEN** `common.NewError(key, err)` 함수가 호출되면
- **THEN** 시스템은 `errors` 맵에 지정된 `key`를 키로, `err.Error()` 문자열을 값으로 설정해야 한다(SHALL)

### Requirement: 비밀번호 해싱

시스템은 사용자 비밀번호를 안전하게 해싱하여 저장해야 한다(MUST).

#### Scenario: 비밀번호 해싱 (setPassword)

- **WHEN** `UserModel.setPassword(password)` 메서드가 호출되면
- **THEN** 시스템은 `bcrypt.GenerateFromPassword`를 `bcrypt.DefaultCost`로 사용하여 비밀번호를 해싱해야 한다(SHALL)
- **THEN** 해싱된 비밀번호를 `UserModel.PasswordHash` 필드에 저장해야 한다(SHALL)

#### Scenario: 빈 비밀번호 거부

- **WHEN** 빈 문자열이 비밀번호로 전달되면
- **THEN** 시스템은 "password should not be empty!" 오류를 반환하고 해싱을 수행하지 않아야 한다(SHALL)

#### Scenario: 비밀번호 검증 (checkPassword)

- **WHEN** `UserModel.checkPassword(password)` 메서드가 호출되면
- **THEN** 시스템은 `bcrypt.CompareHashAndPassword`를 사용하여 입력된 비밀번호와 저장된 해시를 비교해야 한다(SHALL)
- **THEN** 비밀번호가 일치하면 `nil`을, 일치하지 않으면 오류를 반환해야 한다(SHALL)

#### Scenario: 비밀번호 업데이트 시 조건부 해싱

- **WHEN** 사용자 정보 수정 시 비밀번호 필드가 포함되면
- **THEN** 비밀번호 값이 `common.RandomPassword`(더미 값)와 동일한 경우 해싱을 건너뛰어야 한다(SHALL)
- **THEN** 비밀번호 값이 `common.RandomPassword`와 다른 경우에만 새로 해싱하여 저장해야 한다(SHALL)

### Requirement: 유틸리티 함수

시스템은 공통적으로 사용되는 유틸리티 함수를 제공해야 한다(SHALL).

#### Scenario: 랜덤 문자열 생성

- **WHEN** `common.RandString(n)` 함수가 호출되면
- **THEN** 시스템은 `crypto/rand`를 사용하여 영문 대소문자와 숫자로 구성된 `n`자 길이의 암호학적으로 안전한 랜덤 문자열을 생성해야 한다(SHALL)

#### Scenario: 랜덤 정수 생성

- **WHEN** `common.RandInt()` 함수가 호출되면
- **THEN** 시스템은 `crypto/rand`를 사용하여 0 이상 1,000,000 미만의 암호학적으로 안전한 랜덤 정수를 생성해야 한다(SHALL)
