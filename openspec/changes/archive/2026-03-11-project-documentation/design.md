## Context

golang-gin-realworld-example-app은 Go/Gin 기반 RealWorld "Conduit" 백엔드 구현체이다. 현재 코드는 잘 동작하지만, 각 도메인 영역의 요구사항이 코드 안에만 암묵적으로 존재한다. openspec 기반 바이브코딩 워크플로우를 도입했으므로, AI 에이전트가 정확한 구현을 하려면 각 영역의 동작 규격이 명시적 spec 문서로 존재해야 한다.

현재 프로젝트 구조:
- `users/` — 사용자 등록, 로그인, JWT 인증, 프로필, 팔로우
- `articles/` — 아티클 CRUD, 댓글, 즐겨찾기, 태그
- `common/` — DB 초기화, JWT 유틸리티, 유효성 검증, 에러 처리

## Goals / Non-Goals

**Goals:**
- 기존 코드의 동작을 정확히 반영하는 6개 도메인 spec 문서 작성
- 각 spec에 테스트 가능한 시나리오 포함
- README를 바이브코딩 프로젝트에 맞게 업데이트
- 향후 기능 변경 시 spec-first 개발의 기반 마련

**Non-Goals:**
- Go 소스 코드 변경이나 리팩토링
- 새 기능 추가나 API 변경
- 프론트엔드 관련 문서화
- 배포/운영 문서화

## Decisions

### 1. Spec 분류 기준: 도메인 기능 단위
기존 패키지(users/, articles/) 단위가 아닌 도메인 기능 단위로 spec을 분류한다.
- **이유**: `articles/` 패키지가 아티클, 댓글, 즐겨찾기, 태그를 모두 포함하므로 하나의 spec으로는 너무 크다
- **대안**: 패키지 단위 분류 → 거부 (articles 패키지가 4개 도메인을 포함)

### 2. Spec 작성 언어: 한국어
프로젝트 관리자가 한국어를 사용하므로 spec 문서도 한국어로 작성한다.
- **이유**: 리뷰 효율성, 의사소통 명확성
- **대안**: 영어 → RealWorld API spec과 일치하지만 리뷰 비용 증가

### 3. README 업데이트 범위: 최소한
기존 README 구조를 유지하되 바이브코딩 워크플로우 안내만 추가한다.
- **이유**: 상세 내용은 CLAUDE.md, AGENTS.md, openspec/specs/에 이미 있으므로 중복 최소화

## Risks / Trade-offs

- **[Spec과 코드 불일치]** → 기존 코드를 읽고 동작을 정확히 반영하여 작성. 테스트 코드와 교차 검증
- **[Spec 유지보수 부담]** → openspec archive 워크플로우로 변경 시 자동 반영되는 구조 활용
- **[과도한 문서화]** → 각 spec은 RealWorld API spec 기준으로 핵심 요구사항과 시나리오만 포함. 구현 세부사항은 제외
