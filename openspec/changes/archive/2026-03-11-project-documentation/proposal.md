## Why

현재 프로젝트는 기본적인 AGENTS.md와 CLAUDE.md만 있고, 각 모듈(users, articles, common)의 아키텍처, API 스펙, 데이터 모델, 인증 흐름 등에 대한 체계적인 문서가 없다. 바이브코딩 워크플로우에서 AI 에이전트가 정확한 구현을 하려면, 각 도메인 영역의 요구사항과 동작 규격이 openspec specs로 문서화되어야 한다.

## What Changes

- 프로젝트의 핵심 도메인 영역별 openspec spec 문서 생성
- 사용자 인증/프로필 관리 스펙 문서화
- 아티클/댓글/태그 CRUD 스펙 문서화
- 공통 인프라(DB, JWT, 유효성검증) 스펙 문서화
- README를 프로젝트 현황에 맞게 업데이트

## Capabilities

### New Capabilities
- `user-auth`: 사용자 등록, 로그인, JWT 인증, 현재 사용자 조회/수정 API 규격
- `user-profiles`: 프로필 조회, 팔로우/언팔로우 기능 규격
- `articles-crud`: 아티클 생성/조회/수정/삭제, 피드, 목록 필터링/페이지네이션 규격
- `comments`: 아티클 댓글 생성/조회/삭제 규격
- `favorites-and-tags`: 아티클 즐겨찾기, 태그 목록 규격
- `common-infra`: 데이터베이스 설정, JWT 유틸리티, 유효성 검증, 에러 응답 포맷 규격

### Modified Capabilities
<!-- 기존 spec이 없으므로 해당 없음 -->

## Impact

- `openspec/specs/` 하위에 6개의 spec 문서 생성
- README.md 업데이트 (프로젝트 설명, 바이브코딩 워크플로우 안내 추가)
- 기존 Go 코드 변경 없음
- 데이터베이스 마이그레이션 없음
- API 엔드포인트 변경 없음

## Non-goals

- Go 소스 코드 리팩토링이나 기능 변경
- 새로운 API 엔드포인트 추가
- 테스트 코드 수정
- CI/CD 파이프라인 변경
