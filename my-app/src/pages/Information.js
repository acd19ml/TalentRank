import React from 'react';
import { Divider,Typography} from 'antd';
import { BlockMath, InlineMath } from 'react-katex';
import 'katex/dist/katex.min.css';
const { Title, Paragraph, Text, Link } = Typography;


const Formula = () => (
    <Typography>

        <Divider
            style={{
                borderColor: '#7cb305',
            }}
        >
            基础功能
        </Divider>
        <Paragraph>
            <ul  style={{ textAlign: 'left' }}>
                <li>
                    开发者在技术能力方面 TalentRank（类似 Google 搜索的 PageRank），
                    对开发者的技术能力进行评价/评级。评价/评级依据至少包含：项目的重要程度、该开发者在该项目中的贡献度。
                </li>
                <li>
                    开发者的 Nation。有些开发者的 Profile 里面没有写明自己的所属国家/地区。在没有该信息时，可以通过其关系网络猜测其
                    Nation。
                </li>
                <li>
                    开发者的领域。可根据领域搜索匹配，并按 TalentRank 排序。Nation 作为可选的筛选项，比如只需要显示所有位于中国的开发者。
                </li>
                <li>
                    定时获取更新数据库
                </li>
            </ul>
        </Paragraph>
        <Divider
            variant="dashed"
            style={{
                borderColor: '#7cb305',
            }}
            dashed
        >
            高级功能
        </Divider>
        <Paragraph>
            <ul  style={{ textAlign: 'left' }}>
                <li>
                    所有猜测的数据，应该有置信度。置信度低的数据在界面展示为 N/A 值。
                </li>
                <li>
                    开发者技术能力评估信息自动整理。有的开发者在 GitHub
                    上有博客链接，甚至有一些用 GitHub 搭建的网站，也有一些在 GitHub 本身有账号相关介绍。可基于类 ChatGPT 的应用整理出开发者评估信息。
                </li>
            </ul>
        </Paragraph>

        <Divider
            style={{
                borderColor: '#7cb305',
            }}
        >
            贡献度计算公式
        </Divider>
        <div>
            <BlockMath>
                {`- 影响力：\\text{Repository Influence} = \\text{Star} \\times w_{\\text{star}} + \\text{Fork} \\times w_{\\text{fork}} + \\text{Dependents} \\times w_{\\text{dependents}}`}
            </BlockMath>
            <BlockMath>
                {`- 贡献度计算：\\text{Contribution} = \\frac{w_1' \\times \\frac{C_d}{C_t} + w_2' \\times \\frac{I_d}{I_t} + w_3' \\times \\frac{R_d}{R_t} + w_4' \\times \\frac{L_d}{L_t}}{\\sum_{i=1}^{n} w_i'}`}
            </BlockMath>
            <BlockMath>
                {`\\text{Total Score} = \\sum_{\\text{repo}} (\\text{Repository Influence} \\times \\text{Contribution})`}
            </BlockMath>
            <BlockMath>
                {`- 总项目评估：\\text{Overall Score} = \\text{Total Score} \\times (1 + \\text{Followers} \\times w_{\\text{followers}})`}
            </BlockMath>
        </div>
    </Typography>
);

export default Formula;
