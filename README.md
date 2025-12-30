# Commet - Go Gin 웹 애플리케이션

Go Gin 프레임워크 기반의 웹 애플리케이션으로, HTMX와 Alpine.js를 활용한 현대적인 프론트엔드와 PostgreSQL 데이터베이스를 사용합니다.

## 기술 스택

### Backend
- **Go 1.24** - 프로그래밍 언어
- **Gin** - 웹 프레임워크
- **GORM** - ORM
- **JWT** - 인증 (HTTP-Only Cookie)
- **Viper** - 설정 관리

### Frontend
- **HTMX 2.0** - 동적 페이지 업데이트
- **Alpine.js 3.x** - 클라이언트 상태 관리
- **Chart.js 4.x** - 차트 라이브러리
- **Tailwind CSS** - CSS 프레임워크 (CDN)

### Infrastructure
- **PostgreSQL 16** - 데이터베이스
- **Docker Compose** - 컨테이너 오케스트레이션

## 주요 기능

1. **사용자 인증**
   - 회원가입 (이메일, 비밀번호, 이름)
   - 로그인/로그아웃
   - JWT 토큰 기반 세션 (HTTP-Only Cookie)
   - bcrypt 비밀번호 해싱

2. **대시보드**
   - 요약 통계 카드
   - 라인 차트 (월별 매출 추이)
   - 바 차트 (제품별 판매량)
   - 파이 차트 (트래픽 소스)
   - HTMX를 통한 비동기 차트 로딩

## 시작하기

### 사전 요구사항
- Go 1.24+
- Docker & Docker Compose
- Git

### 설치 및 실행

1. **저장소 클론**
```bash
git clone https://github.com/baltop/commet.git
cd commet
```

2. **환경 변수 설정**
```bash
cp .env.example .env
# 필요시 .env 파일 수정
```

3. **PostgreSQL 시작**
```bash
docker compose up -d
```

4. **의존성 설치**
```bash
go mod tidy
```

5. **애플리케이션 빌드 및 실행**
```bash
go build -o bin/server ./cmd/server
./bin/server
```

6. **브라우저에서 접속**
```
http://localhost:8080
```

### 개발 모드 실행
```bash
# air 설치 (선택사항 - 핫 리로드)
go install github.com/air-verse/air@latest

# 개발 서버 실행
air
```

## 프로젝트 구조

```
commet/
├── cmd/
│   └── server/
│       └── main.go              # 애플리케이션 진입점
├── internal/
│   ├── config/                  # 설정 관리
│   ├── database/                # 데이터베이스 연결
│   ├── handlers/                # HTTP 핸들러
│   ├── middleware/              # 미들웨어
│   ├── models/                  # 데이터 모델
│   ├── repository/              # 데이터 액세스
│   └── services/                # 비즈니스 로직
├── web/
│   ├── templates/               # HTML 템플릿
│   └── static/                  # 정적 파일
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

## API 엔드포인트

| Method | Endpoint | 설명 | 인증 |
|--------|----------|------|------|
| GET | / | 홈 (리다이렉트) | - |
| GET | /auth/login | 로그인 페이지 | Guest |
| POST | /auth/login | 로그인 처리 | Guest |
| GET | /auth/register | 회원가입 페이지 | Guest |
| POST | /auth/register | 회원가입 처리 | Guest |
| POST | /auth/logout | 로그아웃 | Auth |
| GET | /dashboard | 대시보드 | Auth |
| GET | /dashboard/charts/line | 라인차트 (HTMX) | Auth |
| GET | /dashboard/charts/bar | 바차트 (HTMX) | Auth |
| GET | /dashboard/charts/pie | 파이차트 (HTMX) | Auth |
| GET | /api/health | 헬스체크 | - |

## 환경 변수

| 변수 | 설명 | 기본값 |
|------|------|--------|
| SERVER_PORT | 서버 포트 | 8080 |
| GIN_MODE | Gin 모드 (debug/release) | debug |
| DB_HOST | DB 호스트 | localhost |
| DB_PORT | DB 포트 | 5435 |
| DB_USER | DB 사용자 | commet |
| DB_PASSWORD | DB 비밀번호 | commet123 |
| DB_NAME | DB 이름 | commet |
| DB_SSLMODE | SSL 모드 | disable |
| JWT_SECRET | JWT 시크릿 키 | - |
| JWT_EXPIRY_HOURS | JWT 만료 시간 | 24 |

## 라이선스

MIT License
