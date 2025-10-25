# 🧊 DeepCool Controller CLI v2
## 📘 Overview

**DeepCool Controller** is a Go-based command-line tool for communicating with DeepCool cooling devices (e.g., **DeepCool LD S360**) via USB HID.  
It continuously monitors CPU metrics—temperature, usage, and power—and sends them to the connected device, allowing dynamic cooling management.  

This version introduces configurable **CPU TDP**, **packet interval**, and optional console output.

## ⚙️ Features

- Real-time CPU monitoring: temperature, usage, and estimated power  
- USB HID device communication  
- Modular architecture (devices, metrics, controllers)  
- CLI interface using [Cobra](https://github.com/spf13/cobra)  
- Optional console output (`--output`)  
- Configurable CPU TDP and packet send interval (`--cpu-tdp-watts`, `--interval`)  
- Easily extensible for new DeepCool devices  

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

### 3. Verify installation

```bash
./deepcool --help
```
## 💻 Usage
### Run the application

```bash
./deepcool run -t <sensorName> -d <deviceModel> [--cpu-tdp-watts <watts>] [--interval <ms>] [--output]
```

**Arguments:**

| Flag                       | Name                                          | Description             | Default | Required |
| -------------------------- | --------------------------------------------- | ----------------------- | ------- | -------- |
| `-t, --temperature-sensor` | Temperature sensor name                       | Example: `k10temp_tctl` | —       | ✅        |
| `-d, --device-model`       | DeepCool device model                         | Example: `dc_ld_s360`   | —       | ✅        |
| `-c, --cpu-tdp-watts`      | Thermal Design Power of CPU in watts          | 170                     | ❌       |          |
| `-i, --interval`           | Interval between packet sends in milliseconds | 500                     | ❌       |          |
| `-o, --output`             | Enable console output of metrics              | false                   | ❌       |          |

**Example:**

```bash
./deepcool run -t k10temp_tctl -d dc_ld_s360 -c 95 -i 300 --output
```

**Example output:**

```
2025/10/25 12:05:21 Device connected!
2025/10/25 12:05:21 Initialization done.
2025/10/25 12:05:21 Sending status: Temp=43.50°C, Usage=27%, Power=26W
2025/10/25 12:05:22 Sending status: Temp=44.12°C, Usage=30%, Power=28W
```
## 🧠 How It Works

1. The CLI initializes the device using **Vendor ID (VID)** and **Product ID (PID)**.
2. The `metrics` package collects CPU temperature, usage, and calculates power based on the configured TDP.
3. The `controllers` package formats the metrics into HID packets and sends them at the configured interval.
4. If the `--output` flag is enabled, metrics are printed to the console.
5. The device dynamically adjusts cooling based on received data.

## 🧰 Commands

| Command            | Description                                  |
| ------------------ | -------------------------------------------- |
| `deepcool run`     | Run the DeepCool monitoring and control loop |
| `deepcool version` | Show the current version of the application  |
| `deepcool help`    | Display help and available commands          |

## 🧩 Adding Support for New Devices

To add a new DeepCool device model:

1. Create a new struct in `devices/` implementing the `Device` interface.
2. Define its **VID**, **PID**, and packet structure (`GetStatusPacket`, `GetDataPacket`, etc.).
3. Add a new `case` in `app.Run()` to register the device model:

```go
case "dc_ld_s280":
    dcDevice = devices.NewDcLdS280()
```

## 🔍 Technical Details

* **Language:** Go 1.22+
* **Libraries Used:**

  * [`spf13/cobra`](https://github.com/spf13/cobra) – CLI framework
  * [`karalabe/hid`](https://github.com/karalabe/hid) – USB HID communication
  * [`shirou/gopsutil`](https://github.com/shirou/gopsutil) – CPU metrics
* **Communication Protocol:** HID (64-byte packets)
* **Configurable Parameters:** CPU TDP, packet interval
* **Supported OS:** Linux / Windows (with HID support)

## 🪪 License

This project is licensed under the MIT License.
See the [LICENSE](LICENSE) file for details.

## 💬 Contributing

1. Fork the repository
2. Create a feature branch:

```bash
git checkout -b feature/new-device
```

3. Commit your changes:

```bash
git commit -am 'Add new DeepCool device support'
```

4. Push and create a Pull Request

---
## 🧊 Example Output
```
2025/10/25 12:10:00 Device connected!
2025/10/25 12:10:00 Initialization done.
2025/10/25 12:10:01 Sending status: Temp=42.85°C, Usage=22%, Power=20W
2025/10/25 12:10:01 Sending status: Temp=43.20°C, Usage=25%, Power=23W
```
**DeepCool Controller CLI v2 — control your cooling, monitor CPU, and stay efficient.**
