import {Group, SimpleCell, Switch} from "@vkontakte/vkui";
import {Icon28MusicOutline} from '@vkontakte/icons';
import {useState} from "react";

export const Tracks = ({tracks}) => {
    let sel = {};
    tracks.forEach((track) => {
        sel = {...sel, [track.id]: false};
    })
    const [selected, setSelected] = useState(sel);

    const toggleSelected = (i) => {
        setSelected({...selected, [i]: !selected[i]})
    }

    return (
        <Group>
            {tracks.map((track) =>
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
