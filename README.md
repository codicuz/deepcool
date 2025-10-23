# 🧊 DeepCool Controller CLI

## 📘 Overview

**DeepCool Controller** is a Go-based command-line interface (CLI) tool that communicates with DeepCool cooling devices (such as **DeepCool LD S360**) via USB HID.  
It monitors real-time CPU metrics (temperature, usage, and power) and sends this data to your DeepCool device for dynamic cooling control.

---

## ⚙️ Features

- 📡 Real-time CPU monitoring: temperature, usage, and estimated power  
- 💡 USB HID device communication  
- 🧩 Modular architecture (devices, metrics, controllers)  
- 🖥️ Command-line interface using [Cobra](https://github.com/spf13/cobra)  
- 🪶 Optional console output for debugging (`--output`)  
- 🧱 Easily extensible for new DeepCool devices  

---

## 🧩 Project Structure

```

deepcool/
├── app/              # Application logic (Run function)
├── cmd/              # CLI command definitions using Cobra
├── controllers/      # HID device controller and communication logic
├── devices/          # Device models (e.g. DeepCool LD S360)
├── metrics/          # CPU metrics collection (temperature, usage, power)
└── main.go           # Entry point

````

---

## 🚀 Installation

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/deepcool.git
cd deepcool
````

### 2. Build the binary

```bash
go build -o deepcool
```

### 3. Verify the installation

```bash
./deepcool --help
```

You should see a list of available commands and flags.

---

## 💻 Usage

### 🔧 Run the application

```bash
./deepcool run -t <sensorName> -d <deviceModel> [--output]
```

**Arguments:**

| Flag                         | Name                              | Description             | Required |
| ---------------------------- | --------------------------------- | ----------------------- | -------- |
| `-t`, `--temperature-sensor` | Temperature sensor name           | Example: `k10temp_tctl` | ✅        |
| `-d`, `--device-model`       | DeepCool device model             | Example: `dc_ld_s360`   | ✅        |
| `-o`, `--output`             | Enable console output for metrics | ❌                       |          |

**Example:**

```bash
./deepcool run -t k10temp_tctl -d dc_ld_s360 --output
```

When `--output` is enabled, you’ll see continuous logs of CPU status:

```
2025/10/24 12:00:11 Device connected!
2025/10/24 12:00:11 Initialization done.
2025/10/24 12:00:12 Sending status: Temp=46.10°C, Usage=32%, Power=54W
```

---

## 🧠 How It Works

1. The CLI initializes the DeepCool device using its **Vendor ID (VID)** and **Product ID (PID)**.
2. The `metrics` package uses `gopsutil` to read system CPU temperature and usage.
3. The `controllers` package formats these values into HID packets.
4. The device receives the packets and updates its cooling behavior dynamically.
5. If `--output` is enabled, metrics are printed to the console in real-time.

---

## 🧰 Commands

| Command            | Description                                  |
| ------------------ | -------------------------------------------- |
| `deepcool run`     | Run the DeepCool monitoring and control loop |
| `deepcool version` | Display current version                      |
| `deepcool help`    | Show help and available commands             |

---

## 🧩 Adding Support for New Devices

To add a new DeepCool device model:

1. Create a new struct in `devices/` implementing the `Device` interface.
2. Define its **VID**, **PID**, and packet structure (`GetStatusPacket`, `GetDataPacket`, etc.).
3. Register the model in `app.Run()` with a new `case` in the device switch block.

Example:

```go
case "dc_ld_s280":
    dcDevice = devices.NewDcLdS280()
```

---

## 🔍 Technical Details

* **Language:** Go 1.22+
* **Libraries:**

  * [`spf13/cobra`](https://github.com/spf13/cobra) – CLI framework
  * [`karalabe/hid`](https://github.com/karalabe/hid) – USB HID communication
  * [`shirou/gopsutil`](https://github.com/shirou/gopsutil) – CPU and system metrics
* **Communication Protocol:** HID (64-byte packets)
* **Supported Systems:** Linux / Windows (with HID support)

---

## 🪪 License

This project is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for more information.

---

## 💬 Contributing

1. Fork this repository
2. Create your feature branch:

   ```bash
   git checkout -b feature/new-device
   ```
3. Commit your changes:

   ```bash
   git commit -am 'Add new DeepCool device support'
   ```
4. Push the branch and open a Pull Request

---

## 🧊 Example Output

```
2025/10/24 15:11:09 Device connected!
2025/10/24 15:11:09 Initialization done.
2025/10/24 15:11:10 Sending status: Temp=42.58°C, Usage=19%, Power=25W
2025/10/24 15:11:11 Sending status: Temp=44.03°C, Usage=23%, Power=30W
```

---

**DeepCool Controller CLI — control your cooling, monitor your CPU, stay cool.**
