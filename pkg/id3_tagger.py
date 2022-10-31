import os
from multiprocessing import Pool, cpu_count
import subprocess

def tag(file):
    artist, title = file.replace('.mp3', '').split(' - ')
    subprocess.call([
            'ffmpeg',
            '-i', file,
            '-metadata', f'artist={artist}',
            '-metadata', f'title={title}',
            '-map', '0:a',
            '-c:a', 'copy',
            f'tagged/{file}'
    ])

if __name__ == '__main__':
    dir_name = os.mkdir('tagged')
    files = os.listdir()
    filtered = [x for x in files if '.mp3' in x]

    cpus_used = cpu_count()

    print(f'Using {cpus_used} cpus')

    with Pool(cpus_used) as p:
        print(p.map(tag, filtered))
