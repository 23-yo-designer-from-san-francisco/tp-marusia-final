from pychorus import find_and_output_chorus
import os
from multiprocessing import Pool, cpu_count

def process(file):
    find_and_output_chorus(file, f'cut_tracks/{file}_chorus', 20)

if __name__ == '__main__':
    dir_name = os.mkdir('cut_tracks')
    files = os.listdir()
    filtered = [x for x in files if '.mp3' in x]

    print(f'Using {cpu_count()} cpus')

    with Pool(cpu_count() - 2) as p:
        print(p.map(process, filtered))
