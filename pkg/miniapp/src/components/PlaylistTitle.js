import {Div, Group, Separator, SimpleCell, Subhead, Title} from "@vkontakte/vkui";
import {Icon28Dice3Outline} from '@vkontakte/icons';
import './PlaylistTitle.css';

export const PlaylistTitle = () => {
    const title = 'Небесная Пыль';
    const generate = () => {
    }
    return (
        <Group>
            <SimpleCell after={<Icon28Dice3Outline/>}>
                Сгенерировать
            </SimpleCell>
            <Separator/>
            <Title className="PlaylistTitle__title">{title.toUpperCase()}</Title>
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
