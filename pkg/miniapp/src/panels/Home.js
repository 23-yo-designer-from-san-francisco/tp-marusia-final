import React, {useEffect, useState} from 'react';
import PropTypes from 'prop-types';

import {Avatar, Button, Group, Header, Panel, PanelHeader, SizeType, Snackbar} from '@vkontakte/vkui';
import {Tracks} from "../components/Tracks";
import {PlaylistTitle} from "../components/PlaylistTitle";
import {BASE_URL} from "../constants";
import {Icon16Done} from "@vkontakte/icons";

const Home = ({id, showLoader}) => {
    const [tracks, setTracks] = useState();
    const [playlistTitle, setPlaylistTitle] = useState('');
    const [selectedTracks, setSelectedTracks] = useState(null);
    const [titleKey, setTitleKey] = useState('');
    const [snackbar, setSnackbar] = useState(null);

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
        setSnackbar(
            <Snackbar
                mode="dark"
                onClose={() => setSnackbar(null)}
                before={
                    <Avatar
                        size={24}
                        style={{background: "var(--vkui--color_background_accent)"}}
                    >
                        <Icon16Done fill="#fff" width={14} height={14}/>
                    </Avatar>
                }
            >
                Плейлист сохранен
            </Snackbar>
        );
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
            <Button class="save-playlist-button" onClick={savePlaylist} sizeY={SizeType.REGULAR}>Сохранить плейлист</Button>
            <Group header={<Header mode="secondary">Доступные треки</Header>}>
                <Tracks updateSelectedTracks={setSelectedTracks} tracks={tracks}/>
            </Group>
            {snackbar}
        </Panel>
    );
};

Home.propTypes = {
    id: PropTypes.string.isRequired,
    go: PropTypes.func.isRequired,
};

export default Home;
