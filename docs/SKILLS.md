# Skills Reference

이 프로젝트에서 사용할 수 있는 Claude Code 스킬(슬래시 커맨드) 가이드입니다.

---

## 개요

스킬은 두 가지 형태로 제공됩니다:

| 유형 | 경로 | 설명 |
|------|------|------|
| **Skills** (`.claude/skills/`) | 자동 트리거 가능 | 설명과 매칭되면 자동으로 호출됨 |
| **Commands** (`.claude/commands/`) | 슬래시 커맨드 전용 | `/opsx:*` 형태로 명시적 호출 |

---

## OpenSpec 워크플로우 스킬

OpenSpec은 이 프로젝트의 계획 중심 개발 워크플로우입니다. 모든 기능 작업은 **탐색 → 제안 → 리뷰 → 구현 → 아카이브** 사이클을 따릅니다.

### `/opsx:explore` — 탐색 모드

> 아이디어를 탐색하고, 문제를 조사하고, 요구사항을 명확히 하는 사고 파트너 모드

- **호출**: `/opsx:explore [주제 또는 change 이름]`
- **목적**: 구현 전 아이디어 탐색, 문제 분석, 옵션 비교
- **특징**:
  - 코드를 읽고 검색할 수 있지만 **코드 작성/구현은 금지**
  - ASCII 다이어그램으로 시각화
  - 기존 OpenSpec change가 있으면 해당 아티팩트를 참조
  - 결론에 도달하면 proposal 생성을 제안
- **입력 예시**:
  - 모호한 아이디어: `/opsx:explore real-time collaboration`
  - 구체적 문제: `/opsx:explore the auth system is getting unwieldy`
  - 옵션 비교: `/opsx:explore postgres vs sqlite for this`
  - change 맥락: `/opsx:explore add-dark-mode`

---

### `/opsx:propose` — 변경 제안

> 새로운 change를 생성하고 모든 아티팩트를 한 번에 만듦

- **호출**: `/opsx:propose <change-name>` 또는 `/opsx:propose <설명>`
- **목적**: 구현할 변경사항의 계획 문서를 자동 생성
- **생성되는 아티팩트**:
  - `proposal.md` — 무엇을 왜 하는지 (What & Why)
  - `design.md` — 어떻게 구현하는지 (How)
  - `tasks.md` — 구현 단계 (Implementation steps)
- **워크플로우**:
  1. change 이름이 없으면 사용자에게 질문
  2. `openspec new change "<name>"` 으로 디렉토리 생성
  3. 아티팩트를 의존성 순서대로 생성
  4. 완료 후 `/opsx:apply`로 구현 시작 안내
- **결과 위치**: `openspec/changes/<change-name>/`

---

### `/opsx:apply` — 변경 구현

> OpenSpec change의 태스크를 순차적으로 구현

- **호출**: `/opsx:apply [change-name]`
- **목적**: 승인된 change의 tasks.md에 정의된 작업을 실행
- **동작 방식**:
  1. change 선택 (이름 지정 또는 자동 추론)
  2. 컨텍스트 파일 읽기 (proposal, design, specs, tasks)
  3. 진행 상황 표시 후 태스크 순차 구현
  4. 각 태스크 완료 시 체크박스 업데이트 (`- [ ]` → `- [x]`)
  5. 모든 태스크 완료 시 archive 제안
- **중단 조건**: 태스크 불명확, 디자인 이슈 발견, 에러/블로커 발생
- **유연성**: 아티팩트가 미완성이어도 태스크가 있으면 실행 가능

---

### `/opsx:archive` — 변경 아카이브

> 완료된 change를 아카이브 디렉토리로 이동

- **호출**: `/opsx:archive [change-name]`
- **목적**: 구현이 완료된 change를 정리하여 보관
- **체크 항목**:
  - 아티팩트 완성도 확인 (미완성 시 경고 후 확인)
  - 태스크 완료 여부 확인 (미완료 시 경고 후 확인)
  - delta spec 동기화 상태 확인
- **결과**: `openspec/changes/<name>` → `openspec/changes/archive/YYYY-MM-DD-<name>/`

---

## Go 개발 스킬

### `go-check` — 테스트 + 린트

> Go 테스트와 린터를 함께 실행하여 코드 품질 검증

- **호출**: `/go-check [패키지 경로]` 또는 자동 트리거
- **실행 내용**:
  - `go test -v -race ./...` — 레이스 컨디션 포함 테스트
  - `golangci-lint run` — 린트 검사
  - 두 명령을 **병렬 실행**
- **결과**: 테스트 통과/실패, 린트 위반 사항 요약
- **활용 시점**: 코드 변경 후, 커밋 전, 구현 검증 시

---

### `go-dev` — 개발 서버 시작

> Go 개발 서버를 시작하고 상태를 확인

- **호출**: `/go-dev [포트]` 또는 자동 트리거
- **동작 순서**:
  1. 기존 서버 프로세스 확인 (포트 충돌 방지)
  2. `go build ./...` 빌드 검증
  3. `go run hello.go &` 백그라운드 실행
  4. health check (`/api/tags` 엔드포인트)
- **기본 URL**: `http://localhost:3000`
- **모델 호출 없음**: 단순 명령 실행만 수행 (빠른 실행)

---

### `go-coverage` — 테스트 커버리지 분석

> 테스트 커버리지 리포트 생성 및 갭 분석

- **호출**: `/go-coverage [패키지 경로]` 또는 자동 트리거
- **실행 내용**:
  - `go test -coverprofile=coverage.out -covermode=atomic ./...`
  - `go tool cover -func=coverage.out`
  - 프로젝트 타겟 대비 비교 (coverage-targets.md 참조)
- **결과**: 패키지별 커버리지 테이블, 갭 함수 목록, 개선 제안
- **타겟 기준**:
  - `articles`: 93%+
  - `users`: 99%+
  - `common`: 85%+

---

## 코드 리뷰 스킬

### `review-code` — 종합 코드 리뷰

> GORM v2 패턴, RealWorld 스펙 준수, 보안을 중심으로 코드 리뷰

- **호출**: `/review-code [파일 또는 패키지]` 또는 자동 트리거
- **리뷰 항목**:
  - **GORM v2 패턴**: Related→Preload, Update→Updates, Delete 포인터, Count 타입, Association 사용법
  - **RealWorld API 스펙**: 상태 코드, 응답 래퍼, 에러 포맷, 페이지네이션
  - **보안**: 인증, 인가, 입력 검증, 비밀번호 보안, DB 보안
  - **Go 모범 사례**: 에러 처리, validator 태그, 미사용 변수
- **범위**: 인수 미지정 시 `git diff`에서 변경된 파일 자동 감지
- **결과**: Critical Issues / Warnings / Suggestions + 최종 판정 (Approved / Changes requested)
- **참조 파일**:
  - `gorm-v2-patterns.md` — GORM 안티패턴 레퍼런스
  - `realworld-api-spec.md` — API 상태 코드 및 응답 포맷
  - `security-checklist.md` — 보안 체크리스트

---

## API 테스트 스킬

### `api-test` — E2E API 테스트

> Newman/Postman을 사용한 RealWorld API 엔드투엔드 테스트

- **호출**: `/api-test [API URL]` 또는 자동 트리거
- **전제 조건**: Newman 설치 (`npm install -g newman`), 서버 실행 중
- **실행 내용**: Postman 컬렉션 기반 API 테스트 실행
- **실패 시**: `realworld-endpoints.md` 참조하여 핸들러/시리얼라이저 분석 후 수정 제안
- **기본 URL**: `http://localhost:3000/api`
- **모델 호출 없음**: 단순 명령 실행만 수행 (빠른 실행)

---

## 워크플로우 다이어그램

```
┌─────────────────────────────────────────────────────────┐
│                   OpenSpec 개발 사이클                     │
└─────────────────────────────────────────────────────────┘

  /opsx:explore          /opsx:propose         /opsx:apply          /opsx:archive
  ┌──────────┐          ┌──────────┐          ┌──────────┐         ┌──────────┐
  │  탐색     │────────▶│  제안     │────────▶│  구현     │───────▶│  아카이브  │
  │          │          │          │          │          │         │          │
  │ 아이디어  │          │ proposal │          │ 태스크    │         │ 완료 정리 │
  │ 문제 분석 │          │ design   │          │ 순차 실행 │         │ 스펙 동기화│
  │ 옵션 비교 │          │ tasks    │          │ 체크박스  │         │ 날짜 태깅 │
  └──────────┘          └──────────┘          └──────────┘         └──────────┘
       │                      │                     │
       │                      │                     │
       ▼                      ▼                     ▼
  ┌──────────┐          ┌──────────┐          ┌──────────┐
  │ go-check │          │ review-  │          │ go-      │
  │ go-dev   │          │ code     │          │ coverage │
  │          │          │          │          │ api-test │
  └──────────┘          └──────────┘          └──────────┘
   개발 & 검증             코드 리뷰             테스트 & 검증
```

---

## 빠른 참조 표

| 스킬 | 커맨드 | 용도 | 자동 트리거 |
|------|--------|------|-------------|
| Explore | `/opsx:explore` | 아이디어 탐색, 문제 조사 | - |
| Propose | `/opsx:propose` | change 생성 + 아티팩트 자동 생성 | - |
| Apply | `/opsx:apply` | 태스크 순차 구현 | - |
| Archive | `/opsx:archive` | 완료된 change 아카이브 | - |
| Go Check | `/go-check` | 테스트 + 린트 병렬 실행 | O |
| Go Dev | `/go-dev` | 개발 서버 시작 | O |
| Go Coverage | `/go-coverage` | 커버리지 리포트 + 갭 분석 | O |
| Review Code | `/review-code` | GORM/RealWorld/보안 종합 리뷰 | O |
| API Test | `/api-test` | E2E API 테스트 (Newman) | O |
