import React, {useEffect, useState} from 'react';
import bridge from '@vkontakte/vk-bridge';
import {
    AdaptivityProvider,
    AppRoot,
    ConfigProvider,
    ScreenSpinner,
    SizeType,
    SplitCol,
    SplitLayout,
    View
} from '@vkontakte/vkui';
import '@vkontakte/vkui/dist/vkui.css';

import Home from './panels/Home';

const App = () => {
    const [scheme, setScheme] = useState('bright_light')
    const [activePanel, setActivePanel] = useState('home');
    const [popout, setPopout] = useState(<ScreenSpinner size='large'/>);

    useEffect(() => {
        bridge.subscribe(({detail: {type, data}}) => {
            if (type === 'VKWebAppUpdateConfig') {
                setScheme(data.scheme)
            }
        });
    }, []);

    const go = e => {
        setActivePanel(e.currentTarget.dataset.to);
    };

    const showLoader = (show) => {
        setPopout(show ? <ScreenSpinner size='large'/> : null)
    }

    return (
        <ConfigProvider scheme={scheme}>
            <AdaptivityProvider sizeX={SizeType.COMPACT}>
                <AppRoot>
                    <SplitLayout popout={popout}>
                        <SplitCol>
                            <View activePanel={activePanel}>
                                <Home showLoader={showLoader} id='home' go={go}/>
                            </View>
                        </SplitCol>
                    </SplitLayout>
                </AppRoot>
            </AdaptivityProvider>
        </ConfigProvider>
    );
}

export default App;
