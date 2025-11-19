# Prometheus Awair Exporter

A Prometheus exporter designed to fetch air quality data from an Awair Element monitor. This project is primarily targeted towards Unraid users.

## How it Works

The exporter connects directly to the Awair Element's local HTTP API, issues a GET request to retrieve sensor data, and then transforms the JSON response into a Prometheus-friendly metrics format.

## Prerequisites

*   A system capable of running Docker containers.
*   A running Prometheus instance.
*   An Awair Element device with its Local API enabled, accessible on the same local network.

## Configuration

This exporter can be configured via a YAML file or environment variables. Environment variables will always override the settings in the YAML file.

### Configuration File

For Unraid users, the recommended method is to create a `config.yml` file in your container's `appdata` directory (which is typically mapped to `/config` inside the container).

**Example `/config/config.yml`:**
```yaml
listen_port: 9101
hosts:
  - 192.168.1.101
  - 192.168.1.102
```

### Environment Variables

| Variable      | Description                                                                 | Default |
|---------------|-----------------------------------------------------------------------------|---------|
| `AWAIR_HOSTS` | A comma-separated list of hostnames or IP addresses for your Awair devices. | `""`      |
| `LISTEN_PORT` | The port on which the exporter will listen for Prometheus scrapes.          | `9101`  |
| `PUID`        | The User ID to run the exporter as.                                         | `99`    |
| `PGID`        | The Group ID to run the exporter as.                                        | `100`   |

## Usage (Docker & Unraid)

To run the container, you need to map the listening port and provide the configuration.

**Example `docker run` command:**
This command uses environment variables for configuration and maps a local directory for persistent configuration, which is common practice on Unraid.

```sh
docker run -d \
  --name awair-exporter \
  -p 9101:9101 \
  -v /mnt/user/appdata/awair-exporter:/config \
  -e PUID=99 \
  -e PGID=100 \
  -e AWAIR_HOSTS="192.168.1.101,192.168.1.102" \
  your-docker-image-name:latest
```

For Unraid, you can use the "Add Container" button in the Docker tab and fill in the appropriate repository, port mappings, volume mappings, and environment variables.

## Prometheus Integration

You will need to add a new scrape configuration to your `prometheus.yml` file to collect metrics from the exporter.

```yaml
scrape_configs:
  - job_name: 'awair'
    static_configs:
      - targets: ['<IP_OF_EXPORTER_HOST>:9101']
```
## Metrics Exposed

The exporter provides the following metrics, all of which are labeled with the `device` host.

| Metric                               | Description                                                      | Units                  |
|--------------------------------------|------------------------------------------------------------------|------------------------|
| `awair_score`                        | Awair's overall air quality score (0-100).                       | -                      |
| `awair_dew_point_celsius`            | The temperature at which air becomes saturated with water vapor. | Celsius                |
| `awair_temperature_celsius`          | The ambient temperature.                                         | Celsius                |
| `awair_humidity_percent`             | The relative humidity.                                           | %                      |
| `awair_absolute_humidity_g_m3`       | The absolute humidity.                                           | g/m³                   |
| `awair_co2_ppm`                      | The measured Carbon Dioxide level.                               | ppm                    |
| `awair_co2_estimated_ppm`            | The estimated Carbon Dioxide level.                              | ppm                    |
| `awair_co2_estimated_baseline`       | The CO2 sensor's internal baseline value.                        | -                      |
| `awair_voc_ppb`                      | The total Volatile Organic Compounds level.                      | ppb                    |
| `awair_voc_baseline`                 | The VOC sensor's internal baseline value.                        | -                      |
| `awair_voc_h2_raw`                   | The raw H2 sensor signal used in the VOC algorithm.              | -                      |
| `awair_voc_ethanol_raw`              | The raw ethanol sensor signal used in the VOC algorithm.         | -                      |
| `awair_pm25_ug_m3`                   | The measured density of particulate matter < 2.5µm.              | µg/m³                  |
| `awair_pm10_estimated_ug_m3`         | The estimated density of particulate matter < 10µm.              | µg/m³                  |

## Next Steps

The project files are complete. The final step is to build the Docker image.

Run the following command from the project's root directory:
```sh
docker build -t awair-exporter:latest .
```

**Note:** This command requires Docker to be installed and running. You may need to run this command with `sudo` or ensure your user is in the `docker` group to avoid permission errors.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

