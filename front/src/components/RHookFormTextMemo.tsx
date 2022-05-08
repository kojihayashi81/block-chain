import { TextField } from "@mui/material"
import { memo } from "react"
import { UseFormReturn } from "react-hook-form"

interface IReactHookFormTextFieldProps {
    methods: UseFormReturn;
    label: string;
    name: string;
    value?: string;
}

const ReactHookFormTextFieldMemo = memo(
    ({ methods, label, name, value }: IReactHookFormTextFieldProps) => (
        <TextField
            label={label}
            variant="outlined"
            error={!!methods.formState.errors[name]}
            value={value}
            helperText={methods.formState.errors[name]?.message ?? ''}
            fullWidth
            margin="dense"
            {...methods.register(name)}
        />
    ),
    (prevProps, nextProps) => {
        return (
            prevProps.methods.formState.isDirty === nextProps.methods.formState.isDirty &&
            prevProps.methods.formState.errors !== nextProps.methods.formState.errors
        );
    }
);

export default ReactHookFormTextFieldMemo;
