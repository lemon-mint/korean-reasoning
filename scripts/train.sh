mkdir train && \
    cd train && \
    sudo apt update && \
    sudo apt install python3-pip python3-venv -y && \
    python3 -m venv .venv

source .venv/bin/activate

pip install wheel setuptools
pip install -U xformers --index-url https://download.pytorch.org/whl/cu124
pip install torch deepspeed tensorboard liger-kernel bitsandbytes accelerate xformers peft trl triton huggingface_hub unsloth hf_transfer lomo-optim
pip install axolotl[flash-attn,deepspeed] --no-build-isolation
#pip install "cut-cross-entropy[transformers] @ git+https://github.com/apple/ml-cross-entropy.git"
export HF_HUB_ENABLE_HF_TRANSFER=1

huggingface-cli login --token $HF_TOKEN