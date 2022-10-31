const tracks = Array.from(document.getElementsByClassName('audio_row')).map((e) => {
    const audioData = JSON.parse(e.dataset.audio)
    const title = audioData[3]
    const artist = audioData[4]
    const audioVkId = `${audioData[1]}_${audioData[0]}`
    const duration = audioData[5]
    return { title, artist, audioVkId, duration }
});

JSON.stringify(tracks);
