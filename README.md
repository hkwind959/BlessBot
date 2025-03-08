# Blockless Bless Network Bot

## Features
- **自动节点交互**
- **多NodeID支持**
- **代理支持**

## Installation

1. Clone the repository to your local machine:
   ```bash
   git clone https://github.com/hkwind959/BlessBot.git
   ```
2. Navigate to the project directory:
   ```bash
   cd BlessBot
   ```
3. Install the necessary dependencies:
   ```bash
   go build
   ```

## Usage
1. Register to blockless bless network account first, if you dont have you can register [https://bless.network/](https://bless.network/dashboard?ref=5PCU0S).
2. Below how to setup this file, put your B7S_AUTH_TOKEN in the text file, example below:
   ```
   eyJhbGcixxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
   ```
   ```bash
   localStorage.getItem('B7S_AUTH_TOKEN')
   ```
3. Set and Modify `config.toml`. Below how to setup this file, example below:
   ```toml
   [[users]]
   user_token = "B7S_AUTH_TOKEN"
   remark = "备注信息"
   nodes = [
      { node_id = "pubKey", proxy = "socks5://username:password@ip:端口", hardware_id = "hardwareId" },
      { node_id = "pubKey", proxy = "socks5://username:password@ip:端口", hardware_id = "hardwareId" },
      { node_id = "pubKey", proxy = "socks5://username:password@ip:端口", hardware_id = "hardwareId" },
      { node_id = "pubKey", proxy = "socks5://username:password@ip:端口", hardware_id = "hardwareId" },
      { node_id = "pubKey", proxy = "socks5://username:password@ip:端口", hardware_id = "hardwareId" },
   ]
   ```