import subprocess
import seaborn as sns
import matplotlib.pyplot as plt
import matplotlib.ticker as ticker
import pandas as pd

tests = ["xsmall", "small", "medium", "large", "xlarge"]
# tests = ["xsmall", "small", "medium"]
cores = [2,4,6,8,12]

seq_data = {"xsmall": [ ], "small": [ ], "medium": [ ], "large": [ ], "xlarge": [ ]}
# seq_data = {"xsmall": [ ], "small": [ ], "medium": [ ]}


for test in tests:
    for i in range(0,5):
        terminalcommand = f"go run benchmark.go s {test}"
        process = subprocess.Popen(terminalcommand, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
        stdout, stderr = process.communicate()
        time = float(stdout.decode().strip())
        seq_data[test].append(time)
        print(f"test: {test}- time: {time}")

for key in seq_data:
    avg = sum(seq_data[key])/len(seq_data[key])
    seq_data[key] = avg


data = {"Test": [ ], "Cores": [ ], "Time": [ ]}

for test in tests:
    for core in cores:
        for i in range(0,5):
            terminalcommand = f"go run benchmark.go p {test} {core}"
            process = subprocess.Popen(terminalcommand, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
            stdout, stderr = process.communicate()
            time = float(stdout.decode().strip())
            
            data["Test"].append(test)
            data["Cores"].append(core)
            speed_up = seq_data[test]/time
            print(f"test: {test} - Cores: {Core} - time: {time}")
            data["Time"].append(speed_up)


df = pd.DataFrame(data)

avg_time = df.groupby(['Test', 'Cores']).mean().reset_index()

for test_type in avg_time['Test'].unique():
    df_type = avg_time[avg_time['Test'] == test_type]
    plt.plot(df_type['Cores'], df_type['Time'], label=test_type)

plt.xlabel('Number of Cores')
plt.ylabel('Speed Up')
plt.title("Average Time per Test as a Function of Cores")
plt.legend()
plt.savefig("speedup1.png")




