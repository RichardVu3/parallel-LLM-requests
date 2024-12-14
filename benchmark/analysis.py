import numpy as np
import pandas as pd
import matplotlib.pyplot as plt

def get_data():
    data = {
        "sequential": 0.0,
        "simple-parallel": {
            "2": 0.0,
            "4": 0.0,
            "6": 0.0,
            "8": 0.0,
            "12": 0.0,
        },
        "work-stealing": {
            "2": 0.0,
            "4": 0.0,
            "6": 0.0,
            "8": 0.0,
            "12": 0.0,
        }
    }

    with open("benchmark/runtime/sequential.txt", "r") as f:
        times = f.readlines()
    data["sequential"] = np.mean([float(time.strip()) for time in times])
    for thread in ["2", "4", "6", "8", "12"]:
        with open(f"benchmark/runtime/simple-parallel/{thread}.txt", "r") as f:
            times = f.readlines()
        data["simple-parallel"][thread] = np.mean([float(time.strip()) for time in times])
        with open(f"benchmark/runtime/work-stealing/{thread}.txt", "r") as f:
            times = f.readlines()
        data["work-stealing"][thread] = np.mean([float(time.strip()) for time in times])

    return data

def get_speedup(data):
    speedup = {"simple-parallel": {}, "work-stealing": {}}
    for thread in ["2", "4", "6", "8", "12"]:
        speedup["simple-parallel"][thread] = data["sequential"] / data["simple-parallel"][thread]
        speedup["work-stealing"][thread] = data["sequential"] / data["work-stealing"][thread]
    return speedup

def graph_speedup(speedup):
    threads = ["2", "4", "6", "8", "12"]
    plt.plot(threads, [speedup["simple-parallel"][thread] for thread in threads], label="simple-parallel")
    plt.plot(threads, [speedup["work-stealing"][thread] for thread in threads], label="work-stealing")
    plt.ylim(1.75)
    plt.xlabel("Number of threads")
    plt.ylabel("Speedup")
    plt.legend()
    plt.title("Speedup of simple-parallel and work-stealing")
    plt.savefig("benchmark/my-speedup.png")
    print("Graph saved as benchmark/my-speedup.png")

if __name__ == "__main__":
    data = get_data()
    speedup = get_speedup(data)
    graph_speedup(speedup)