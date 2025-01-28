# Datasets

## open-thoughts-114k.jsonl

https://huggingface.co/datasets/lemon-mint/OpenThoughts-114k-JSONL/blob/main/data-train.jsonl

- Python 코드 생성 과제
- 수학 추론 과제 & LaTeX 출력
- 과학 QA

```html
"<|begin_of_thought|>
(추론과정)
\n\n<|end_of_thought|>
\n\n<|begin_of_solution|>\n\n
(최종응답)
\n\n<|end_of_solution|>"
```

1. 질문에서 공통된 접두어 제거
2. 응답 특수 토큰 기준으로 분리 & 특수 토큰 존재 여부 확인
3. 특수 토큰 기준으로 응답 분리
4. 처리된 데이터 셋 "ot-processed-01.jsonl" 저장