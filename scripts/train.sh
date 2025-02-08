mkdir train && \
    cd train && \
    sudo apt update && \
    sudo apt install python3-pip python3-venv -y && \
    python3 -m venv .venv

source .venv/bin/activate

pip install wheel setuptools ninja packaging
pip install transformers liger-kernel torch wandb bitsandbytes accelerate peft trl huggingface_hub hf_transfer lomo-optim
pip install flash-attn --no-build-isolation
#pip install "cut-cross-entropy[transformers] @ git+https://github.com/apple/ml-cross-entropy.git"
export HF_HUB_ENABLE_HF_TRANSFER=1
export WANDB_PROJECT=phi-4-ko-reasoning-v1

huggingface-cli login --token $HF_TOKEN