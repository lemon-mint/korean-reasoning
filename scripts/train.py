from transformers import AutoModelForCausalLM, AutoTokenizer, TrainingArguments
AutoLigerKernelForCausalLM
from trl import SFTTrainer
from cut_cross_entropy.transformers import cce_patch

BASE_MODEL = "sh2orc/Llama-3.1-Korean-8B-Instruct"
TOKENIZER_MODEL = "unsloth/Meta-Llama-3.1-8B-Instruct"

print("Loading the base model")
model = AutoModelForCausalLM.from_pretrained(BASE_MODEL)
model = cce_patch(model, reduction="none")

tokenizer = AutoTokenizer.from_pretrained(TOKENIZER_MODEL)

def formatting_prompts_func(examples):
    convos = examples["messages"]
    texts = [tokenizer.apply_chat_template(convo, tokenize = False, add_generation_prompt = False) for convo in convos]
    return { "text" : texts, }
pass

print("Loading the dataset")
from datasets import load_dataset
dataset = load_dataset("lemon-mint/korean-reasoning-v02", split = "train")
dataset = dataset.map(formatting_prompts_func, batched = True,)

training_args = TrainingArguments(
    output_dir="output",
    per_device_train_batch_size=8,

)
