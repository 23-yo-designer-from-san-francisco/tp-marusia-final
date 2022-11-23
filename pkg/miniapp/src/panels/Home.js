import React from 'react';
import PropTypes from 'prop-types';
import {tracks} from '../tracks';

import {Group, Header, Panel, PanelHeader} from '@vkontakte/vkui';
import {Tracks} from "../components/Tracks";
import {PlaylistTitle} from "../components/PlaylistTitle";

const Home = ({id, go}) => (
    <Panel id={id}>
        <PanelHeader>Создать плейлист</PanelHeader>
        <Group header={<Header mode="secondary">Название плейлиста</Header>}>
            <PlaylistTitle/>
        </Group>
        <Group header={<Header mode="secondary">Доступные треки</Header>}>
            <Tracks tracks={tracks}/>
        </Group>
    </Panel>
);

Home.propTypes = {
    id: PropTypes.string.isRequired,
    go: PropTypes.func.isRequired,
};

export default Home;
