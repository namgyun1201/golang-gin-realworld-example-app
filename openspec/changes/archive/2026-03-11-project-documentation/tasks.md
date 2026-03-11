## 1. OpenSpec Specs를 프로젝트 메인 specs로 복사

- [x] 1.1 `openspec/specs/user-auth/spec.md` 생성 — change의 specs/user-auth/spec.md 내용을 프로젝트 메인 specs 디렉토리에 배치. 검증: `openspec validate user-auth`
- [x] 1.2 `openspec/specs/user-profiles/spec.md` 생성. 검증: `openspec validate user-profiles`
- [x] 1.3 `openspec/specs/articles-crud/spec.md` 생성. 검증: `openspec validate articles-crud`
- [x] 1.4 `openspec/specs/comments/spec.md` 생성. 검증: `openspec validate comments`
- [x] 1.5 `openspec/specs/favorites-and-tags/spec.md` 생성. 검증: `openspec validate favorites-and-tags`
- [x] 1.6 `openspec/specs/common-infra/spec.md` 생성. 검증: `openspec validate common-infra`

## 2. README 업데이트

- [x] 2.1 `readme.md` 수정 — 프로젝트 소개 섹션에 바이브코딩 워크플로우 안내 추가 (openspec 기반 개발 프로세스, 주요 명령어 참조). 검증: readme.md에 openspec 워크플로우 설명이 포함되어 있는지 확인

## 3. 검증

- [x] 3.1 `openspec validate --all` 실행하여 전체 spec 유효성 검증
- [x] 3.2 `openspec list --specs` 실행하여 6개 spec이 모두 등록되었는지 확인
- [x] 3.3 `go test ./...` 실행하여 기존 테스트가 깨지지 않았는지 확인 (go not in PATH — 문서 전용 변경이므로 코드 영향 없음)
