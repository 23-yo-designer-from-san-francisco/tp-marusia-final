const map = {};
Array.from(document.getElementsByClassName('audio_row')).forEach((e) => {
    const audioData = JSON.parse(e.dataset.audio)
    const title = audioData[3]
    const artist = audioData[4]
    const audioVkId = `${audioData[1]}_${audioData[0]}`
    const duration = audioData[5]
    map[`${title}${artist}`] = { ...map[`${title}${artist}`], title, artist }
    map[`${title}${artist}`][`duration_${duration}`] = audioVkId;
});
const tracks = [];
for (const [,value] of Object.entries(map)) {
    tracks.push(value);
}
tracks.sort((a, b) => a.title < b.title ? -1 : a.title > b.title ? 1 : 0);

JSON.stringify(tracks);
