import {Div, Group, Separator, SimpleCell, Subhead, Title} from "@vkontakte/vkui";
import {Icon28Dice3Outline} from '@vkontakte/icons';
import './PlaylistTitle.css';
import {BASE_URL} from "../constants";
import {useEffect, useState} from "react";

export const PlaylistTitle = ({title, showLoader, updateTitleKey}) => {
    const [playlistTitle, setPlaylistTitle] = useState('');
    const [titleKey, setTitleKey] = useState('');

    useEffect(() => {
        load();
    }, []);

    async function load() {
        showLoader(true);
        const response = await fetch(`${BASE_URL}/playlists/generate-title`);
        const {title_key, title} = await response.json();
        setPlaylistTitle(title);
        setTitleKey(title_key);
        updateTitleKey(title_key);
        showLoader(false);
    }

    return (
        <Group>
            <SimpleCell onClick={load} after={<Icon28Dice3Outline/>}>
                Сгенерировать
            </SimpleCell>
            <Separator/>
            <Title className="PlaylistTitle__title">{playlistTitle.toUpperCase()}</Title>
            <Div className="PlaylistTitle__subtitle-container">
                <Subhead className="PlaylistTitle__subtitle">
                    Это название вашего плейлиста. Вы можете поделиться этим
                    плейлистом с друзьями. Для выбора
                    плейлиста в скиле скажите "Плейлист"
                </Subhead>
            </Div>
        </Group>
    );
}
