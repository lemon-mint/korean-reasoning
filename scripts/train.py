from transformers import AutoModelForCausalLM, AutoTokenizer, TrainingArguments
#from liger_kernel.transformers import AutoLigerKernelForCausalLM
#from cut_cross_entropy.transformers import cce_patch
from trl import SFTTrainer, SFTConfig
import torch

BASE_MODEL = "google/gemma-2-9b-it"
TOKENIZER_MODEL = "google/gemma-2-9b-it"
OUTPUT_MODEL = "exp-models/gemma-2-ko-reasoning-v1"
DATASET = "exp-models/korean-reasoning-mixture-20250203-preview"

print("Loading the base model")
model = AutoModelForCausalLM.from_pretrained(
    BASE_MODEL,
    torch_dtype=torch.bfloat16,
    attn_implementation="flash_attention_2",
)
#model = cce_patch(model)
tokenizer = AutoTokenizer.from_pretrained(TOKENIZER_MODEL)

def formatting_prompts_func(examples):
    convos = examples["messages"]
    texts = [tokenizer.apply_chat_template(convo, tokenize = False, add_generation_prompt = False) for convo in convos]
    return { "text" : texts, }
pass

print("Loading the dataset")
from datasets import load_dataset
dataset = load_dataset(DATASET, split = "train")
dataset = dataset.map(formatting_prompts_func, batched = True,)

output_dir = "outputs"
max_seq_length = 8192

training_args = SFTConfig(
    num_train_epochs=2,
    #max_steps=5,
    max_seq_length = max_seq_length,
    packing = True, # Can make training 5x faster for short sequences.

    per_device_train_batch_size=8,
    #gradient_accumulation_steps=1,
    gradient_checkpointing=True,

    learning_rate = 4e-5,
    lr_scheduler_type="cosine",
    warmup_ratio=0.05,

    bf16=True,

    use_liger_kernel=True,
    use_liger=True,

    weight_decay=0.01,
    optim="adalomo",

    report_to="wandb",
    output_dir=output_dir,
    push_to_hub=True,
    hub_model_id=OUTPUT_MODEL,
    hub_strategy="checkpoint",
    save_steps=300,
    save_total_limit=1,
    logging_steps=1,
)

trainer = SFTTrainer(
    model = model,
    processing_class=tokenizer,
    train_dataset = dataset,
    args = training_args,
)

print("do train")
trainer.train()
print("end train")

model.push_to_hub(OUTPUT_MODEL+"-iter1")
tokenizer.push_to_hub(OUTPUT_MODEL+"-iter1")

model.save_pretrained(OUTPUT_MODEL+"-iter1")
tokenizer.save_pretrained(OUTPUT_MODEL+"-iter1")