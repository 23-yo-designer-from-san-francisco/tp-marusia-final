import {FormItem, FormLayout, Group, IconButton, Input, SimpleCell, Switch} from "@vkontakte/vkui";
import {Icon16Clear, Icon28MusicOutline} from '@vkontakte/icons';
import React, {useEffect, useState} from "react";

export const Tracks = ({tracks, updateSelectedTracks}) => {
    if (!tracks) {
        return null;
    }
    let sel = {};
    tracks.forEach((track) => {
        sel = {...sel, [track.id]: false};
    })
    const [selected, setSelected] = useState(sel);
    const [renderedTracks, setRenderedTracks] = useState(tracks);
    const textInput = React.createRef();
    const clear = () => {
        textInput.current.value = "";
        setRenderedTracks(tracks);
    };

    useEffect(() => {
        updateSelectedTracks(selected);
    }, [selected]);

    const toggleSelected = (i) => {
        setSelected({...selected, [i]: !selected[i]})
    }

    const search = () => setRenderedTracks(tracks.filter((track) => {
        const trackFullTitle = `${track.artist} ${track.title}`.toLowerCase();
        return trackFullTitle.includes(textInput.current.value.toLowerCase())
            || trackFullTitle.split(' ').reverse().join(' ').includes(textInput.current.value.toLowerCase())
    }));

    return (
        <Group>
            <FormLayout>
                <FormItem top="Поиск">
                    <Input
                        getRef={textInput}
                        type="text"
                        placeholder="Начните вводить исполнителя или название песни"
                        onChange={search}
                        after={
                            <IconButton
                                hoverMode="opacity"
                                aria-label="Очистить поле"
                                onClick={clear}
                            >
                                <Icon16Clear/>
                            </IconButton>
                        }
                    />
                </FormItem>
            </FormLayout>
            {renderedTracks.map((track) =>
                <SimpleCell
                    key={track.id}
                    before={<Icon28MusicOutline/>}
                    after={
                        <Switch readOnly checked={selected[track.id]}/>
                    }
                    onClick={() => toggleSelected(track.id)}
                >
                    {`${track.artist} — ${track.title}`}
                </SimpleCell>
            )}
        </Group>
    )
}
