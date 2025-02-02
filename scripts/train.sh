mkdir train && \
    cd train && \
    sudo apt update && \
    sudo apt upgrade -y && \
    sudo apt install python3-pip python3-venv -y && \
    python3 -m venv .venv

source .venv/bin/activate

pip install wheel setuptools ninja packaging
pip install -U xformers --index-url https://download.pytorch.org/whl/cu124
pip install torch wandb liger-kernel bitsandbytes accelerate peft trl triton huggingface_hub hf_transfer lomo-optim
pip install flash-attn --no-build-isolation
pip install "cut-cross-entropy @ git+https://github.com/apple/ml-cross-entropy.git"
export HF_HUB_ENABLE_HF_TRANSFER=1
export WANDB_PROJECT=LLaMa-3.1-Korean-Reasoning-8B-Instruct

huggingface-cli login --token $HF_TOKEN