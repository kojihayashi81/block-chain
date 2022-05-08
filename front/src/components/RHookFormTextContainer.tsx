import { FC } from "react"
import { useFormContext } from "react-hook-form"

import ReactHookFormTextFieldMemo from "./RHookFormTextMemo"

interface ReactHookFormTextFieldContainerProps {
    name: string;
    label: string;
    value?: string;
}

const ReactHookFormTextFieldContainer: FC<ReactHookFormTextFieldContainerProps> = ({
    name,
    label,
    value,
}: ReactHookFormTextFieldContainerProps) => {
    const methods = useFormContext();

    return <ReactHookFormTextFieldMemo methods={methods} label={label} name={name} value={value} />;
};

export default ReactHookFormTextFieldContainer;
