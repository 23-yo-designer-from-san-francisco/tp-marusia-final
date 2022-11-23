import React, {useEffect, useState} from 'react';
import PropTypes from 'prop-types';

import {Button, Group, Header, Panel, PanelHeader, SizeType} from '@vkontakte/vkui';
import {Tracks} from "../components/Tracks";
import {PlaylistTitle} from "../components/PlaylistTitle";
import {BASE_URL} from "../constants";

const Home = ({id, showLoader}) => {
    const [tracks, setTracks] = useState();
    const [playlistTitle, setPlaylistTitle] = useState('');
    const [selectedTracks, setSelectedTracks] = useState(null);
    const [titleKey, setTitleKey] = useState('');

    async function fetchData() {
        showLoader(true);
        const response = await fetch(`${BASE_URL}/music`);
        const jsonResponse = await response.json();
        console.log(jsonResponse);
        setTracks(jsonResponse);
        showLoader(false);
    }

    async function savePlaylist() {
        const tr = [];
        console.log(1, selectedTracks);
        for (const [key, value] of Object.entries(selectedTracks)) {
            if (value) {
                tr.push(tracks[key - 1].id)
            }
        }
        await fetch(`${BASE_URL}/playlists/create`, {
            method: 'POST',
            body: JSON.stringify({
                title_key: titleKey,
                music_ids: tr,
            })
        });
    }

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <Panel id={id}>
            <PanelHeader>Создать плейлист</PanelHeader>
            <Group header={<Header mode="secondary">Название плейлиста</Header>}>
                <PlaylistTitle updateTitleKey={setTitleKey} showLoader={showLoader} title={playlistTitle}/>
            </Group>
            <Button onClick={savePlaylist} sizeY={SizeType.REGULAR}>Сохранить плейлист</Button>
            <Group header={<Header mode="secondary">Доступные треки</Header>}>
                <Tracks updateSelectedTracks={setSelectedTracks} tracks={tracks}/>
            </Group>
        </Panel>
    );
};

Home.propTypes = {
    id: PropTypes.string.isRequired,
    go: PropTypes.func.isRequired,
};

export default Home;
